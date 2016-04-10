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

function get(path, success, error) {
	var xhttp = new XMLHttpRequest();
	xhttp.responseType = "json";
	xhttp.onreadystatechange = function() {
		if (xhttp.readyState == 4) {
			if (Math.round(xhttp.status/100) == 2) {
				return success(xhttp.response);
			}

			error && error(xhttp.response);
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
	showAbout: function() {
		var about = document.getElementById("dialog-about");
		about.showModal();
	},
	getInitialState: function() {
		return {
			version: "Loading...",
			packs: []
		};
	},
	updateStatus: function(ev) {
		var packs = ev.detail.Packs;
		var installed = packs.Installed.map(
			function(p) {
				return {
					name: p.Name,
					installed: !p.Installing,
					count: p.NameCount
				}
			}
		);
		var available = packs.Available.map(
			function(p) {
				return {
					name: p.Name,
					file: p.File,
					installed: false,
					count: p.NameCount
				}
			}
		);

		this.setState({
			version: ev.detail.Version,
			packs: installed.concat(available)
		});
	},
	componentDidMount: function() {
		document.addEventListener("Status", this.updateStatus.bind(this));
		get("/status", function(resp) {
			publish("Status", resp);
		});
	},
	componentWillUnmount: function() {
		document.removeEventListener("Status", this.updateStatus.bind(this));
	},
	showPacks: function() {
		var packs = document.getElementById("dialog-packs");
		packs.showModal();
	},
	render: function() {
		return (
			el("div", {id: "menu"},
				 el("button", {id:"menu-button", className:"mdl-button mdl-js-button mdl-js-ripple-effect mdl-button--icon"},
						el("i", {className:"material-icons"}, "more_vert")
				 ),
				 el("ul", {className:"mdl-menu mdl-menu--bottom-left mdl-js-menu mdl-js-ripple-effect", htmlFor:"menu-button"},
						el("li", {className: "mdl-menu__item"}, "Help"),
						el("li", {className: "mdl-menu__item", onClick: this.showPacks.bind(this)}, "Packs"),
						el("li", {className: "mdl-menu__item"}, "Settings"),
						el("li", {className: "mdl-menu__item", onClick: this.showAbout.bind(this)}, "About")
				 ),
				 el(Packs, {packs: this.state.packs}),
				 el(About, {version: this.state.version})
			)
		);
	}
});

var Packs = React.createClass({
	displayName: "Packs",
	close: function() {
		var packs = document.getElementById("dialog-packs");
		packs.close();
	},
	componentDidUpdate: function() {
		componentHandler.upgradeDom();
	},
	render: function() {
		var packs = [];
		for (var i = 0; i < this.props.packs.length; i++) {
			var p = this.props.packs[i];
			p.idx = i;
			packs.push(el(Pack, p));
		}

		return (
			el("dialog", {className: "mdl-dialog", id: "dialog-packs"},
				 el("h4", {className: "mdl-dialog__title"}, "Packs"),
				 el("div", {className: "mdl-dialog__content"},
						el("div", {},
							 el("ul", {className: "mdl-list"}, packs)
						)
				 ),
				 el("div", {className: "mdl-dialog__actions"},
						el("button", {className: "mdl-button close", type: "button", onClick:this.close.bind(this)}, "Dismiss")
				 )
			)
		);
	}
});

var warn = function(msg) {
	alert(msg);
};

var Pack = React.createClass({
	displayName: "Pack",
	install: function() {
		get(
			"/install/"+this.props.file,
			function () {this.setState({installed: true});}.bind(this),
			function () { warn("Failed to install package.");}.bind(this)
		);
	},
	remove: function () {
		get(
			"/remove/"+this.props.name,
			function () {
				console.log("setting installed: false for", this);
				document.getElementById("pack-checkbox-label-"+this.props.idx).MaterialCheckbox.uncheck();
				this.setState({installed: false});
			}.bind(this),
			function () {
				console.log("remove failed for some reason", arguments);
			}.bind(this)
		);
	},
	toggleInstall: function() {
		console.log("Pack props", this.props);
		if (this.state.installed) {
			return this.remove();
		}
		return this.install();
	},
	getInitialState: function () {
		return {
			installed: this.props.installed
		};
	},
	render: function() {
		var checkBoxId = "pack-checkbox-"+this.props.idx;
		var checkBoxLabelId = "pack-checkbox-label-"+this.props.idx;
		var checkBoxOptions = {
			className: "mdl-checkbox__input",
			type:"checkbox",
			id:checkBoxId,
			onClick: this.toggleInstall.bind(this)
		};
		if (this.state.installed) {
			checkBoxOptions.checked = "checked";
		}
		return (
			el("li", {className: "mdl-list__item mdl-list__item--two-line"},
				 el("span", {className: "mdl-list__item-primary-content"},
						el("i", {className: "material-icons mdl-list__item-avatar"}, "info"),
						el("span", {className: "mdl-list__item-text"}, this.props.name),
						el("span", {className: "mdl-list__item-sub-title"}, this.props.count + " entries")
				 ),
				 el("span", {className: "mdl-list__item-secondary-content"},
						el("label", {id: checkBoxLabelId, className: "mdl-list__item-secondary-action mdl-checkbox mdl-js-checkbox mdl-js-ripple-effect", htmlFor:checkBoxId},
							 el("input", checkBoxOptions)
						)
				 )
			)
		);
	}
});

var About = React.createClass({
	displayName: "About",
	close: function() {
		var about = document.getElementById("dialog-about");
		about.close();
	},
	render: function() {
		return (
			el("dialog", {className: "mdl-dialog", id: "dialog-about"},
				 el("h4", {className: "mdl-dialog__title"}, "About"),
				 el("div", {className: "mdl-dialog__content"},
						el("p", {},
							 "More information at ",
							 el("a", {href:"https://github.com/fgeller/rad"}, "github.com/fgeller/rad")
						),
						el("p", {}, "Build version: "+this.props.version)
				 ),
				 el("div", {className: "mdl-dialog__actions"},
						el("button", {className: "mdl-button close", type: "button", onClick:this.close.bind(this)}, "Dismiss")
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

var loadDoc = function (target) {
	var ifrm = document.getElementById("ifrm");
	if (ifrm.src != target) {
		ifrm.src = target;
	}
}

var loadDocDelayed = (function() {
	var timer = 0;
	var delay = 1000;
	return function(target) {
		clearTimeout(timer);
		timer = setTimeout(function () { loadDoc(target); }, delay);
	};
})();

var SearchResult = React.createClass({
	displayName: "SearchResult",
	componentDidMount: function() {
		if (this.props.selected) {
			loadDocDelayed(this.props.target);
		}
	},
	componentDidUpdate: function() {
		componentHandler.upgradeDom();
	},
	select: function (e) {
		loadDoc(this.props.target);
		publish("SelectSearchResult", this.props.index);
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
		this.setState(ns);
	},
	moveSelectionUp: function() {
		var n = Math.max(0, this.state.selection - 1);
		publish("SelectSearchResult", n);
	},
	moveSelectionDown: function() {
		var n = Math.min(this.state.results.length - 1, this.state.selection + 1);
		publish("SelectSearchResult", n);
	},
	componentDidMount: function () {
		document.addEventListener("SearchResults", this.updateResults.bind(this));
		document.addEventListener("SelectSearchResult", this.updateSelection.bind(this));

		key.unbind('down');
		key('down', function(event) {
			this.moveSelectionDown();
			event.cancelBubble = true;
			return false;
		}.bind(this));

		key.unbind('up');
		key('up', function(event) {
			this.moveSelectionUp();
			event.cancelBubble = true;
			return false;
		}.bind(this));
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
