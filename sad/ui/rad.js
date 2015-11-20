'use strict';

var SettingsDefaults = {
	resultLimit: 5,
	requestThrottle: 100,
	autoLoad: false
};

// http://stackoverflow.com/a/2880929/750284
function urlParams() {
	var urlParams = {};
	var match,
		pl     = /\+/g,  // Regex for replacing addition symbol with a space
		search = /([^&=]+)=?([^&]*)/g,
		decode = function (s) { return decodeURIComponent(s.replace(pl, " ")); },
		query  = window.location.search.substring(1);

	while (match = search.exec(query)) {
		urlParams[decode(match[1])] = decode(match[2]);
	}

	return urlParams;
}

function addHistory(name, value) {
	var ps = urlParams();
	ps[name] = value;

	var entry = "/"
	var first = true;
	for (var k in ps) {
		entry += first ? "?" : "&";
		first = false;
		entry +=  encodeURIComponent(k) + "=" + encodeURIComponent(ps[k])
	}

	window.history.replaceState(null, "", entry);
}

// http://stackoverflow.com/a/105074
function guid() {
	function s4() {
		return Math.floor((1 + Math.random()) * 0x10000)
			.toString(16)
			.substring(1);
	}
	return (
		s4() + s4() + '-' +
			s4() + '-' +
			s4() + '-' +
			s4() + '-' +
			s4() + s4() + s4()
	);
}

function get(path, success) {
	var xhttp = new XMLHttpRequest();
	xhttp.responseType = "json";
	xhttp.onreadystatechange = function() {
		if (xhttp.readyState == 4 && xhttp.status == 200) {
			if (xhttp.status == 200) {
				return success(xhttp.response);
			}

			console.log("Failed to react to non-200 xhttp:", xhttp);
		}

	};
	xhttp.open("GET", path, true);
	xhttp.send();
}

function socket() {
	var host = window.location.hostname;
	var port = window.location.port;
	var conn = new WebSocket("ws://"+host+":"+port+"/s");
	conn.onopen = function () { console.log("WebSocket open.") }
	conn.onerror = function (err) { console.log("WebSocket error", err) }
	conn.onmessage = function (msg) { console.log("Got message", msg) }
	conn.onclose = function () { console.log("Connection closed", arguments); }
	return conn;
}

key('/', function(event) {
	console.log("setting search-field focus");
	document.getElementById("search-field").focus();
	event.cancelBubble = true;
	return false;
});

key.filter = function(event){
	var tagName = (event.target || event.srcElement).tagName;
	if (event.target && event.target.id && event.target.id == "search-field") {
		return true;
	}
	return !(tagName == 'INPUT' || tagName == 'SELECT' || tagName == 'TEXTAREA');
}

var publish = function(name, data) {
	var ev = new CustomEvent(name, {detail: data});
	console.debug(name, data);
	document.dispatchEvent(ev);
}

//
// Components
//

var el = React.createElement;

var Menu = React.createClass({
	displayName: "Menu",
	render: function() {
		return (
			el("div", {id: "menu"},
				el("button", {id:"menu-button", className:"mdl-button mdl-js-button mdl-js-ripple-effect mdl-button--icon"},
					el("i", {className:"material-icons"}, "more_vert")
				),
				el("ul", {className:"mdl-menu mdl-menu--bottom-left mdl-js-menu mdl-js-ripple-effect", htmlFor:"menu-button"},
					el("li", {className: "mdl-menu__item"}, "Help"),
					el("li", {className: "mdl-menu__item"}, "Packages"),
					el("li", {className: "mdl-menu__item"}, "Settings")
				)
			)
		);
	}
});

var SearchButton = React.createClass({
	displayName: "SearchButton",
	render: function() {
		return (
			el("div", {},
				el("label", {className: "mdl-button mdl-js-button mdl-button--icon", htmlFor: "search-field"},
					el("i", {className: "material-icons"}, "search")
				)
			)
		);
	}
});

var SearchField = React.createClass({
	displayName: "SearchField",
	getInitialState: function() {
		return {};
	},
	parseQuery: function(q) {
		var ps = { Limit: 5, Pack: "", Path: "", Member: "" };
		var qs = q.split(" ");
		if (qs.length == 1) {
			ps.Pack = qs[0];
		} else if (qs.length == 2) {
			ps.Pack = qs[0];
			ps.Member = qs[1];
		} else if (qs.length > 2) { // TODO: maybe join 2++
			ps.Pack = qs[0];
			ps.Path = qs[1];
			ps.Member = qs[2];
		}
		return ps;
	},
	highlightLabels: function(ps) {
		var pck = document.getElementById("search-label-pack");
		pck.className = (ps.Pack.length > 0) ? "highlight" : "";
		var pth = document.getElementById("search-label-path");
		pth.className = (ps.Path.length > 0) ? "highlight" : "";
		var mem = document.getElementById("search-label-member");
		mem.className = (ps.Member.length > 0) ? "highlight" : "";
	},
	search: function() {
		var query = this.refs.searchFieldInput.getDOMNode().value;
		var ps = this.parseQuery(query);
		this.highlightLabels(ps);
		this.streamSearch(ps);
	},
	streamSearch: function(ps) {
		if (this.state.sock) {
			this.state.sock.close();
		}

		var sock = socket();
		this.setState({sock: sock, results: []});
		publish("SearchResults", []);

		sock.onmessage = function(msg) {
			var entry = JSON.parse(msg.data);
			this.setState({results: this.state.results.concat([entry])});
			publish("SearchResults", this.state.results);
		}.bind(this);

		sock.onopen = function() {
			sock.send(JSON.stringify(ps));
		}.bind(this);

		sock.onclose = function() {
			console.log("Finished request, found " + this.state.results.length + " results, params:", ps);
		}.bind(this);
	},
	render: function() {
		return (
			el("div", {className:"mdl-textfield mdl-js-textfield mdl-textfield--floating-label"},
				el("input", {
					className: "mdl-textfield__input",
					ref: "searchFieldInput",
					type: "text",
					id: "search-field",
					onChange: this.search.bind(this)
				}),
				el("label", {className:"mdl-textfield__label", htmlFor: "search-field"},
					el("span", {id: "search-label-pack"}, "pack"),
					" ",
					el("span", {id: "search-label-path"}, "path"),
					" ",
					el("span", {id: "search-label-member"}, "member")
				)
			)
		);
	}
});

var SearchBar = React.createClass({
	displayName: "SearchBar",
	render: function() {
		return (
			el("div", {id: "search-fields"},
				el(Menu, {}),
				el(SearchField, {}),
				el(SearchButton, {})
			)
		);
	}
});

var SearchResult = React.createClass({
	displayName: "SearchResult",
	componentDidUpdate: function() {
		componentHandler.upgradeDom();
	},
	select: function (e) {
		this.loadDocumentation();
		publish("SelectSearchResult", this.props.index);
	},
	loadDocumentation: function () {
		document.getElementById("ifrm").src = this.props.target;
	},
	render: function() {
		var cn = this.props.selected ? "is-selected" : "";
		return (
			el("tr", {id:"search-result-"+this.props.index, className: cn, onClick: this.select.bind(this)},
				el("td", {className: "mdl-data-table__cell--non-numeric"}, this.props.member),
				el("td", {className: "mdl-data-table__cell--non-numeric"}, this.props.path)
			)
		);
	}
});

var SearchResults = React.createClass({
	displayName: "SearchResults",
	getInitialState: function() {
		return {selection: 0, results: []};
	},
	updateResults: function(ev) {
		this.setState({results: ev.detail, selection: 0});
	},
	updateSelection: function(ev) {
		var n = ev.detail;
		var ns = {selection: n};
		console.log("update selection", ns);
		this.setState(ns);
	},
	componentDidMount: function () {
		document.addEventListener("SearchResults", this.updateResults.bind(this));
		document.addEventListener("SelectSearchResult", this.updateSelection.bind(this));
	},
	componentDidUpdate: function() {
		componentHandler.upgradeDom();
	},
	componentWillUnmount: function () {
		document.removeEventListener("SearchResults", this.updateResults.bind(this));
		document.removeEventListener("SelectSearchResult", this.updateSelection.bind(this));
	},
	instantiateSearchResult: function(index, data) {
		return (
			el(SearchResult, {
				index: index,
				selected: this.state.selection === index,
				member: data["Member"],
				path: data["Namespace"],
				target: data["Target"]
			})
		);
	},
	moveSelectionLeft: function() {
		var n = Math.max(0, this.state.selection - 1);
		publish("SelectSearchResult", n);
	},
	moveSelectionRight: function() {
		var n = Math.min(this.state.results.length - 1, this.state.selection + 1);
		publish("SelectSearchResult", n);
	},
	domResults: function() {
		var results = [];
		for (var i = 0; i < this.state.results.length; i++) {
			var r = this.state.results[i];
			results.push(this.instantiateSearchResult(i, r));
		}
		return results;
	},
	render: function() {
		var d = this.state.results.length == 0 ? "none" : "block";
		return (
			el("div", {id: "search-results", style: {display: d}},
				el("table", {className:"mdl-data-table mdl-js-data-table"},
					el("tbody", {}, this.domResults())
				)
			)
		);
	}
});

var Search = React.createClass({
	displayName: "Search",
	render: function() {
		return (
			el("div", {id: "search"},
				el(SearchBar, {}),
				el(SearchResults, {})
			)
		);
	}
});

var DocumentationFrame = React.createClass({
	displayName: "DocumentationFrame",
	render: function() {
		return (
			el("div", {id: "doc-container"},
				el("iframe", {id: "ifrm"})
			)
		);
	}
});

var Rad = React.createClass({
	displayName: "Rad",
	componentDidUpdate: function() {
		componentHandler.upgradeDom();
	},
	render: function() {
		return (
			el("div", {},
				el(Search, {}),
				el(DocumentationFrame, {})
			)
		);
	}
});

ReactDOM.render(
	el(Rad, {}),
	document.getElementById("rad")
);
