'use strict';

var SettingsDefaults = {
	resultLimit: 10,
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

//
// Components
//

var Menu = React.createClass({
	render: function() {
		return (
			<div id="menu">
				<button id="menu-button" className="mdl-button mdl-js-button mdl-js-ripple-effect mdl-button--icon">
					<i className="material-icons">more_vert</i>
				</button>
				<ul className="mdl-menu mdl-menu--bottom-left mdl-js-menu mdl-js-ripple-effect" htmlFor="menu-button">
					<li className="mdl-menu__item">Help</li>
					<li className="mdl-menu__item">Packages</li>
					<li className="mdl-menu__item">Settings</li>
				</ul>
			</div>
		);
	}
});

var Search = React.createClass({
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
	getInitialState: function() {
		return {};
	},
	streamSearch: function(ps) {
		if (this.state.sock) {
			this.state.sock.close();
		}
		var sock = socket();
		this.setState({sock: sock, results: []});

		sock.onmessage = function(msg) {
			var entry = JSON.parse(msg.data);
			this.setState({results: this.state.results.concat([entry])});
			console.log("results: ", this.state.results);
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
			<div id="search">
				<div>
					<label id="search-button" className="mdl-button mdl-js-button mdl-js-ripple-effect mdl-button--icon" htmlFor="search-field-input">
						<i className="material-icons">search</i>
					</label>
				</div>
				<div>
					<form action="#">
						<div className="mdl-textfield mdl-js-textfield mdl-textfield--floating-label">
							<input
								id="search-field-input"
								ref="searchFieldInput"
								className="mdl-textfield__input"
								type="text"
								onChange={this.search}
							/>
							<label className="mdl-textfield__label" htmlFor="search-field-input">
								<span id="search-label-pack">pack</span>&nbsp;
								<span id="search-label-path">path</span>&nbsp;
								<span id="search-label-member">member</span>
							</label>
						</div>
					</form>
				</div>
			</div>
		);
	}
});

var SearchResults = React.createClass({
	render: function() {
		return (
			<div className="mdl-tabs mdl-js-tabs mdl-js-ripple-effect">
				<div className="mdl-tabs__tab-bar">
					<a className="mdl-tabs__tab scrollindicator">
						<i className="material-icons scrollindicator scrollindicator--left disabled"></i>
					</a>
					<a className="mdl-tabs__tab">
						<div className="member">*core-java-api*</div>
						<div className="path">clojure.java.javadoc</div>
					</a>
					<a className="mdl-tabs__tab">
						<div className="member">andThen</div>
						<div className="path">scala.concurrent.FutureFutureFuture</div>
					</a>
					<a className="mdl-tabs__tab">
						<div className="member">columnNumberColumnNumber</div>
						<div className="path">Error.prototypeprototypeprototype</div>
					</a>
					<a className="mdl-tabs__tab scrollindicator">
						<i className="material-icons scrollindicator scrollindicator--right"></i>
					</a>
				</div>
			</div>
		);
	}
});

var Rad = React.createClass({
	componentDidUpdate: function() {
		componentHandler.upgradeDom();
	},
	render: function() {
		return (
			<div id="nav">
				<div className="mdl-grid mdl-grid--no-spacing">
					<div className="mdl-cell mdl-cell--5-col">
						<Menu />
						<Search />
					</div>
					<div className="mdl-cell mdl-cell--7-col">
						<SearchResults />
					</div>
				</div>
			</div>
		);
	}
});

React.render(
	<Rad />,
	document.getElementById("rad")
);
