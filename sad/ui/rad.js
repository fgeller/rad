'use strict';

var SettingsDefaults = {
	resultLimit: 10,
	requestThrottle: 100,
	autoLoad: false
};

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

var InfoButton = React.createClass({
	render: function() {
		return (
			<div
				id={this.props.id}
				className="info-button"
				style={{ display: this.props.display }}
				onClick={this.props.toggleInfo} >
				<i className="fa fa-info-circle"></i>
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
	loadHelp: function () {
		document.getElementById("ifrm").src = "/readme.html"
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
				<i
					id="search-icon"
					className="fa fa-search"
				></i>
				<i
					id="search-help"
					className="fa fa-question"
					onClick={this.loadHelp}
				></i>
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
		if (memName.length > 16) {
			memName = memName.substring(0, 16) + "...";
		}
		var clsName = "search-result"
		if (this.props.selected) {
			clsName += " selected-search-result";
			if (this.props.autoLoad) {
				this.open();
			}
		}

		if (this.props.visible) {
			clsName += " visible-search-result";
		}

		return (
			<div className={clsName} onClick={this.open}>
				<div className="search-result-label">{this.props.label}</div>
				<div className="member-name">{memName}</div>
				<div className="namespace">{namespace}</div>
			</div>
		);
	}
});

var Pack = React.createClass({
	getInitialState: function() {
		return {
			installing: this.props.installing || false,
			showDescription: false,
			removing: false
		};
	},
	remove: function(e) {
		e.stopPropagation();
		get(
			"/remove/"+this.props.name,
			function(xhttp) { this.props.loadPacks(); }.bind(this)
		);
		this.props.loadPacks();
	},
	install: function(e) {
		e.stopPropagation();
		get(
			"/install/"+this.props.file,
			function(xhttp) { this.props.loadPacks(); }.bind(this)
		);
		this.props.loadPacks();
	},
	toggleInfo: function(e) {
		this.setState({showDescription: !this.state.showDescription});
	},
	render: function() {
		var created = (new Date(this.props.created)).toISOString().substring(0, 10);
		var displayRemove = this.props.installed ? "inline-block" : "none";
		var displayInstall = this.props.installed ? "none" : "inline-block";
		var desc = this.props.description;
		var hasDesc = desc && desc.length > 0;
		var infoRowStyle = { display: this.state.showDescription ? 'inline-block' : 'none' };
		return (
			<div className="settings-pack">
				<div className="settings-pack-row">
					<div className="settings-pack-row-value">{this.props.name}
					</div><div className="settings-pack-row-value">{this.props.type}
					</div><div className="settings-pack-row-value">{this.props.version}
					</div><div className="settings-pack-row-value">{created}
					</div><div
									className="settings-pack-row-value settings-pack-button" >
						<InfoButton
							toggleInfo={ this.toggleInfo }
							display={ hasDesc ? 'inline-block' : 'none' }
						/>
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
				<div className="settings-pack-row settings-pack-desc" style={infoRowStyle} dangerouslySetInnerHTML={{__html: desc}}>
				</div>
			</div>
		);
	}
});

	var RequestSettings = React.createClass({
		getInitialState: function () {
			return {
				resultLimit: this.props.resultLimit || 3,
				requestThrottle: this.props.requestThrottle || 100,
				autoLoad: this.props.autoLoad || false,
				resultLimitInvalid: false,
				requestThrottleInvalid: false
			}
		},
		update: function () {
			var newResultLimit = document.getElementById("settings-result-limit").value;
			if (!newResultLimit || isNaN(Number(newResultLimit)) || Number(newResultLimit) <= 0 || Number(newResultLimit) >= 500) {
				newResultLimit = SettingsDefaults.resultLimit;
			}
			var newRequestThrottle = document.getElementById("settings-request-throttle").value;
			if (!newRequestThrottle || isNaN(Number(newRequestThrottle)) || Number(newRequestThrottle) <= 0) {
				newRequestThrottle = SettingsDefaults.requestThrottle;
			}
			var newAutoLoad = document.getElementById("settings-auto-load").checked;
			localStorage.clear();
			localStorage.setItem("SettingsResultLimit", newResultLimit);
			localStorage.setItem("SettingsRequestThrottle", newRequestThrottle);
			localStorage.setItem("SettingsAutoLoad", newAutoLoad);
			console.log("Updated settings in local storage: ", localStorage);
			this.props.updateSettings();
		},
		componentWillReceiveProps: function(nextProps) {
			var rl = document.getElementById("settings-result-limit").value;
			var rt = document.getElementById("settings-request-throttle").value;
			this.setState({
				resultLimitInvalid: (rl != nextProps.resultLimit),
				resultLimitInvalidValue: rl,
				requestThrottleInvalid: (rt != nextProps.requestThrottle),
				requestThrottleInvalidValue: rt
			});
		},
		render: function () {

			var rlValue = this.state.resultLimitInvalid ? this.state.resultLimitInvalidValue : this.props.resultLimit
			var rtValue = this.state.requestThrottleInvalid ? this.state.requestThrottleInvalidValue : this.props.requestThrottle;
			return (
				<div id="settings-general">
					<div className="settings-header">Settings</div>
					<div className="settings-label-value">
						<div className="settings-label">Result limit</div>
						<div className="settings-value">
							<input
								id="settings-result-limit"
								type="text"
								value={rlValue}
								className={this.state.resultLimitInvalid ? "input-invalid" : ""}
								onChange={this.update}
							/>
						</div>
						<div className="settings-explanation">Limit the number of results that are requested.</div>
					</div>
					<div className="settings-label-value">
						<div className="settings-label">Request throttle</div>
						<div className="settings-value">
							<input
								id="settings-request-throttle"
								type="text"
								value={rtValue}
								className={this.state.requestThrottleInvalid ? "input-invalid" : ""}
								onChange={this.update}
							/>
						</div>
						<div className="settings-explanation">Send search queries at most every x milliseconds.</div>
					</div><div className="settings-label-value">
						<div className="settings-label">Load immediately</div>
						<div className="settings-value">
							<input
								id="settings-auto-load"
								type="checkbox"
								checked={this.props.autoLoad}
								onChange={this.update}
							/>
						</div>
						<div className="settings-explanation">Automatically load current selection in document frame.</div>
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
					this.props.updatePacks(data.Packs.Installed)
				}.bind(this)
			);
		},
		componentWillMount: function() {
			this.loadPacks();
		},
		componentWillUpdate: function() {
			if (_.some(this.state.availablePacks, function (p) { return p.Installing })) {
				setTimeout(function() { this.loadPacks(); }.bind(this), 5000);
			}
		},
		hide: function() {
			document.getElementById("settings-container").style.visibility = "hidden";
		},
		stopEvent: function(e) {
			e.stopPropagation();
		},
		genUID: function() {
			return ''+Math.round(Math.random()*10000000000)
		},
		render: function() {
			var installed = _.sortBy(this.state.installedPacks, function(p) { return p.Name; }).map(function(p, idx) {
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
						description={p.Description}
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
						installing={p.Installing}
						description={p.Description}
						installed={false} />
				);
			}.bind(this)
			);

			var installedDom = <div />;
			if (installed.length > 0) {
				installedDom = <div id="settings-installed-packs">
						<div className="settings-header">Installed Packs</div>
						<div className="settings-pack-row">
							<div className="settings-pack-row-label">Name</div>
							<div className="settings-pack-row-label">Type</div>
							<div className="settings-pack-row-label">Version</div>
							<div className="settings-pack-row-label">Created</div>
						</div>
						{ installed }
				</div>;
			}

			var availableDom = <div />;
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
					<div id="settings-content" onClick={this.stopEvent}>
						{ installedDom }
						{ availableDom }
						<RequestSettings
							resultLimit={this.props.resultLimit}
							requestThrottle={this.props.requestThrottle}
							autoLoad={this.props.autoLoad}
							updateSettings={this.props.updateSettings}
						/>
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
				results: [],
				settings: this.readSettings()
			};
		},
		installThrottledSearch: function (throttle) {
			var throttledSearch = _.debounce(
				function (txt) { this.streamSearch(txt); },
				throttle,
				{ leading: false, trailing: true }
			);
			this.throttledSearch = throttledSearch;
		},
		updatePacks: function (installed) {
			if (installed.length == 1 && this.state.query.length == 0) {
				this.setState({query: installed[0].Name + ' '});
			}
		},
		updateSettings: function () {
			var newSettings = this.readSettings();
			this.installThrottledSearch(newSettings.requestThrottle);
			this.setState({settings: newSettings});
		},
		readSettings: function() {
			var result = {
				requestThrottle: Number(localStorage["SettingsRequestThrottle"] || SettingsDefaults.requestThrottle),
				resultLimit: Number(localStorage["SettingsResultLimit"] || SettingsDefaults.resultLimit),
				autoLoad: localStorage["SettingsAutoLoad"] === "true"
			};
			return result;
		},
		search: function(text) {
			this.setState({query: text, selected: 0, results: []});
			this.throttledSearch(text);
		},
		streamSearch: function(text){
			window.history.replaceState(null, "", "/?q="+encodeURIComponent(text));
			var qs = text.split(" ");
			if (qs.length < 2) {
				return;
			}
			this.props.sock.close();
			var pk = qs[0];
			var pt = qs[1];
			var m	= qs[2] || "";
			var lim = this.state.settings.resultLimit;
			var req = {"Limit": lim, "Pack": pk, "Path": pt, "Member": m};

			this.props.sock = socket();
			this.props.sock.onmessage = function(msg) {
				var entry = JSON.parse(msg.data);
				this.setState({results: this.state.results.concat([entry])});
			}.bind(this);
			this.props.sock.onopen = function() {
				this.props.sock.send(JSON.stringify(req));
			}.bind(this);
			this.props.sock.onclose = function() {
				console.log("Finished request [" + text + "] found " + this.state.results.length + " results.");
				if (this.state.results.length == 0) {
					key.unbind('return');
				}
				if (this.state.results.length == 0 && this.state.settings.autoLoad) {
					document.getElementById("ifrm").src = "/zero.html";
				}
			}.bind(this);
		},
		selectResult: function(idx) {
			if (idx >= 0 && idx < this.state.results.length) {
				this.setState({selected: idx});
			}
		},
		shiftSelection: function(left, right) {
			// if left: try -1
			var sub1 = this.state.selected - 1
			if (left) {
				if (sub1 >= 0) { this.setState({selected: sub1}) }
				return
			}
			// if right: try +1
			var add1 = this.state.selected + 1
			if (right) {
				if (add1 < this.state.results.length) { this.setState({selected: add1}) }
				return
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

			document.getElementById("search-field").focus();

			this.installThrottledSearch(this.state.settings.requestThrottle);

			var params = window.location.search.substring(1);
			var arrParam = params.split("=");

			if (this.state.query == "" &&
				arrParam.length == 2 &&
				arrParam[0] == "q" &&
				arrParam[1].length) {
					var q = decodeURIComponent(arrParam[1].replace(/\+/g,'%20'));
					this.search(q);
			}
		},
		showSettings: function () {
			document.getElementById("settings-container").style.visibility = "visible";
		},
		render: function(){
			var entries = [];
			for (var i = 0; i < this.state.results.length; i++) {
				var entry = this.state.results[i];
				var isVisible = (
					(this.state.selected >= this.state.results.length-2 &&
						i >= this.state.results.length-3) ||
						(i >= this.state.selected-1 &&
							((this.state.selected > 0 && i <= this.state.selected+1) ||
								(this.state.selected == 0 && i <= 2)))
				);
				var label = ""+(i+1)+"/"+this.state.results.length;
				if (this.state.results.length == this.state.settings.resultLimit) {
					label += "+";
				}
				entries.push(
					<SearchResult
						entry={entry}
						index={i}
						label={label}
						visible={isVisible}
						autoLoad={this.state.settings.autoLoad}
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
					<Settings
						resultLimit={this.state.settings.resultLimit}
						requestThrottle={this.state.settings.requestThrottle}
						autoLoad={this.state.settings.autoLoad}
						updateSettings={this.updateSettings}
						updatePacks={this.updatePacks}
					/>
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
