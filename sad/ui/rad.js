'use strict';

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

var Loading = React.createClass({
	render: function() {
		return (
			<div
				id={this.props.id}
				className="loading-screen"
				style={{ display: this.props.display }} >
				<div className="loading-spinner">
					<i className="fa fa-spinner fa-pulse"></i>
				</div>
			</div>
		);
	}
});

var InstallButton = React.createClass({
	render: function() {
		return (
			<div
				id={this.props.id}
				className="install-button"
				style={{ display: this.props.display }}
				onClick={this.props.install} >
				<i className="fa fa-cloud-download"></i>
			</div>
		);
	}
});

var RemoveButton = React.createClass({
	render: function() {
		return (
			<div
				id={this.props.id}
				className="remove-button"
				style={{ display: this.props.display }}
				onClick={ this.props.remove } >
				<i className="fa fa-times"></i>
			</div>
		);
	}
});

var SearchField = React.createClass({
	search: function() {
		var query = this.refs.search.getDOMNode().value;
		this.props.search(query)
	},
	render: function() {
		return (
			<div>
				<input
					id="search-field"
					type="text"
					ref="search"
					placeholder="Search here..."
					value={this.props.query}
					onChange={this.search} />
				<i id="search-icon" className="fa fa-search"></i>
			</div>
		);
	}
});

var SearchResult = React.createClass({
	open: function() {
		if (!this.props.selected){
			this.props.selectResult(this.props.index);
		}

		var href = this.props.entry["Target"];
		console.log("Opening target", href);
		document.getElementById("ifrm").src = href;
	},
	render: function() {
		if (this.props.selected) {
			key.unbind('return');
			key('return', function(event) {
				this.open();
				event.cancelBubble = true;
				return false;
			}.bind(this));
		}

		// \u00a0 = &nbsp;
		var namespace = this.props.entry["Namespace"] || "\u00a0";
		var memName = this.props.entry["Member"] || "\u00a0";
		if (memName.length > 20) {
			memName = memName.substring(0, 20) + "...";
		}
		var clsName = "search-result"
		if (this.props.selected) {
			clsName += " selected-search-result";
		}
		if (this.props.index == 0) {
			clsName += " first-search-result";
		}
		return (
			<div className={clsName} onClick={this.open}>
				<div className="member-name">{memName}</div>
				<div className="namespace">{namespace}</div>
			</div>
		);
	}
});

var Pack = React.createClass({
	getInitialState: function() {
		return {
			installing: false,
			removing: false
		};
	},
	remove: function(e) {
		e.stopPropagation();
		this.setState({ removing: true });
		get(
			"/remove/"+this.props.name,
			function(xhttp) {
				this.setState({removing: false});
				this.props.loadPacks();
			}.bind(this)
		);
	},
	install: function(e) {
		e.stopPropagation();
		this.setState({ installing: true });
		get(
			"/install/"+this.props.file,
			function(xhttp) {
				this.setState({installing: false});
				this.props.loadPacks();
			}.bind(this)
		);
	},
	render: function() {
		var created = (new Date(this.props.created)).toISOString().substring(0, 10);
		var displayRemove = this.props.installed ? "inline-block" : "none";
		var displayInstall = this.props.installed ? "none" : "inline-block";
		return (
			<div className="settings-pack">
				<div className="settings-pack-row">
					<div className="settings-pack-row-value">{this.props.name}
					</div><div className="settings-pack-row-value">{this.props.type}
					</div><div className="settings-pack-row-value">{this.props.version}
					</div><div className="settings-pack-row-value">{created}
					</div><div
									className="settings-pack-row-value settings-pack-button"
									style={{display: displayInstall}} >
						<InstallButton
							install={ this.install }
							display={ this.state.installing ? "none" : "inline-block" } />
						<Loading
							display={ this.state.installing ? "inline-block" : "none" } />
					</div><div
									className="settings-pack-row-value settings-pack-button"
									style={{display: displayRemove}}>
						<RemoveButton
							remove={ this.remove }
							display={ this.state.removing ? "none" : "inline-block" } />
						<Loading
							display={ this.state.removing ? "inline-block" : "none" } />
					</div>
				</div>
			</div>
		);
	}
});

var Settings = React.createClass({
	getInitialState: function() {
		return {
			installedPacks: [],
			availablePacks: []
		};
	},
	loadPacks: function() {
		get(
			"/status",
			function(data) {
				this.setState({
					installedPacks: data.Packs.Installed,
					availablePacks: data.Packs.Available,
				});
			}.bind(this)
		);
	},
	componentWillMount: function() {
		this.loadPacks();
	},
	hide: function() {
		document.getElementById("settings-container").style.visibility = "hidden";
	},
	genUID: function() {
		return ''+Math.round(Math.random()*10000000000)
	},
	render: function() {
		var installed = this.state.installedPacks.map(function(p, idx) {
			var id = guid();
			return (
				<Pack
					name={p.Name}
					uid={id}
					key={id}
					loadPacks={this.loadPacks}
					type={p.Type}
					version={p.Version}
					created={p.Created}
					installed={true} />
			);
		}.bind(this));

		var available = this.state.availablePacks.filter(
			function(p) {
				var installed = false;
				this.state.installedPacks.forEach(function(ip) {
					if (p.Name == ip.Name && p.Version == ip.Version && p.Created == ip.Created) {
						installed = true;
					}
				});
				return !installed;
			}.bind(this)
		).map(function(p, idx) {
				var id = guid();
				return (
					<Pack name={p.Name}
						uid={id}
						key={id}
						type={p.Type}
						version={p.Version}
						created={p.Created}
						file={p.File}
						loadPacks={this.loadPacks}
						installed={false} />
				);
			}.bind(this)
		);

		var availableDom = <div />; // might want to add more here
		if (available.length > 0) {
			availableDom = <div id="settings-available-packs">
						<div className="settings-header">Available Packs</div>
						<div className="settings-pack-row">
							<div className="settings-pack-row-label">Name</div>
							<div className="settings-pack-row-label">Type</div>
							<div className="settings-pack-row-label">Version</div>
							<div className="settings-pack-row-label">Created</div>
						</div>
						{ available }
			</div>;
		}

		return (
			<div id="settings-container" onClick={this.hide}>
				<div id="settings-content">
					<div id="settings-installed-packs">
						<div className="settings-header">Installed Packs</div>
						<div className="settings-pack-row">
							<div className="settings-pack-row-label">Name</div>
							<div className="settings-pack-row-label">Type</div>
							<div className="settings-pack-row-label">Version</div>
							<div className="settings-pack-row-label">Created</div>
						</div>
						{ installed }
					</div>
					{ availableDom }
				</div>
			</div>
		);
	}
});

var Search = React.createClass({
	getInitialState: function(){
		return {
			selected: 0,
			query:'',
			results: []
		};
	},
	search: function(text) {
		this.setState({query: text, selected: 0, results: []});
		this.streamSearch(text);
	},
	streamSearch: function(text){
		this.setState({query: text, selected: 0, results: []});
		var qs = text.split(" ");
		if (qs.length < 2) {
			return;
		}
		this.props.sock.close();
		var pk = qs[0]
		var pt = qs[1]
		var m	= qs[2] || ""
		var req = {
			"Limit": 3,
			"Pack": pk,
			"Path": pt,
			"Member": m
		};

		this.props.sock = socket();
		this.props.sock.onmessage = function(msg) {
			var entry = JSON.parse(msg.data);
			this.setState({results: this.state.results.concat([entry])});
		}.bind(this);
		this.props.sock.onopen = function() {
			this.props.sock.send(JSON.stringify(req));
		}.bind(this);
		this.props.sock.onclose = function() {
			console.log("Finished request [" + text + "].");
		}.bind(this);
	},
	selectResult: function(idx) {
		if (idx >= 0 && idx <= 3 && idx < this.state.results.length) {
			this.setState({selected: idx});
		}
	},
	shiftSelection: function(left, right) {
		// if left: try -1
		var sub1 = this.state.selected - 1
		if (left && sub1 >= 0) {
			this.setState({selected: sub1});
			return;
		}
		// if right: try +1
		var add1 = this.state.selected + 1
		if (right && add1 < 3 && add1 < this.state.results.length) { // TODO magic number
			this.setState({selected: add1});
			return;
		}
		// if none: 0
		if (!left && !right) {
			this.setState({selected: 0});
			return;
		}
	},
	componentDidMount: function() {
		key.unbind('shift+right');
		key.unbind('shift+left');

		key('shift+right', function(event) {
			this.shiftSelection(false, true);
			event.cancelBubble = true;
			return false;
		}.bind(this));

		key('shift+left', function(event) {
			this.shiftSelection(true, false);
			event.cancelBubble = true;
			return false;
		}.bind(this));
	},
	showSettings: function () {
		document.getElementById("settings-container").style.visibility = "visible";
	},
	render: function(){
		var entries = [];
		for (var i = 0; i < this.state.results.length && i < 3; i++) {
			var entry = this.state.results[i];
			entries.push(
				<SearchResult
					entry={entry}
					index={i}
					selected={i == this.state.selected}
					selectResult={this.selectResult} />
			);
		}

		var params = window.location.search.substring(1);
		var arrParam = params.split("=");

		var doc = "/readme.html";
		if (arrParam.length == 2 && arrParam[0] == "doc") {
			doc = arrParam[1] + decodeURIComponent(window.location.hash).replace(/&gt;/g, ">").replace(/&lt;/g, "<");
		}

		return (
			<div id="main-container">
				<Settings />
				<div id="menu-container">
					<i className="fa fa-cogs" onClick={this.showSettings}></i>
				</div>
				<div id="search-field-container">
					<SearchField query={this.state.query} search={this.search}/>
				</div>
				<div id="search-result-container">
					{entries}
				</div>
				<div id="ifrm-container">
					<iframe id="ifrm" src={doc} />
				</div>
			</div>
		);
	}
});

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

React.render(
	<Search sock={socket()} />,
	document.getElementById("search")
);
