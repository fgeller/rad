var Loading = React.createClass({
    render: function() {
	return <div id={this.props.id} className="loading-screen"><div className="loading-spinner"><i className="fa fa-spinner fa-pulse"></i></div><div className="loading-message">{ this.props.message }</div></div>
    }
});

var InstallButton = React.createClass({
    render: function() {
	return <div id={this.props.id} className="install-button" onClick={this.props.install}><i className="fa fa-cloud-download"></i></div>
    }
});

var SearchField = React.createClass({
    search: function() {
        var query = this.refs.search.getDOMNode().value;
        this.props.search(query)
    },
    render: function() {
        return <div>
                 <input
                      id="search-field"
                      type="text"
                      ref="search"
                      placeholder="Search here..."
                      value={this.props.query}
                      onChange={this.search} />
                 <i id="search-icon" className="fa fa-search"></i>
               </div>
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

        // \00a0 = &nbsp;
        var namespace = this.props.entry["Namespace"] || "\00a0";
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
        return <div className={clsName} onClick={this.open}>
                 <div className="member-name">{memName}</div>
                 <div className="namespace">{namespace}</div>
               </div>
    }
});

var Pack = React.createClass({
    install: function(e) {
	var btn = $("#install-button-"+this.props.id+".install-button");
	var load = $("#loading-screen-"+this.props.id+".loading-screen");
	btn.css('display', 'none');
	load.css('display', 'block');

        $.get(
            "/install/"+this.props.file,
            {},
            function(data, flag) {
		this.props.loadPacks();
            }.bind(this)
        );
	return e.stopPropagation();
    },
    render: function() {
	var created = (new Date(this.props.created)).toISOString().substring(0, 10);
        return <div className="settings-pack">
                 <div className="settings-pack-row">
                     <div className="settings-pack-row-value">{this.props.name}
                     </div><div className="settings-pack-row-value">{this.props.type}
                     </div><div className="settings-pack-row-value">{this.props.version}
                     </div><div className="settings-pack-row-value">{created}
                     </div><div className="settings-pack-row-value" style={{display: this.props.installed ? "none" : "inline-block"}}>
                     <InstallButton id={"install-button-"+this.props.id} install={this.install} />
                     <Loading id={"loading-screen-"+this.props.id} message="" />
                   </div>
                 </div>
               </div>
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
        $.get(
            "/status/packs/installed",
            {},
            function(data, flag) {
                var installedPacks = [];
                _.forEach(data, function(p) { installedPacks.push(p) });
                this.setState({installedPacks: installedPacks});
            }.bind(this),
            "json" // we expect json
        );
        $.get(
            "/status/packs/available",
            {},
            function(data, flag) {
                var availablePacks = [];
                _.forEach(data, function(p) { availablePacks.push(p) });
                this.setState({availablePacks: availablePacks});
            }.bind(this),
            "json" // we expect json
        );
    },
    componentWillMount: function() {
	this.loadPacks();
    },
    hide: function() {
        $("#settings-container").css("visibility", "hidden");
    },
    render: function() {
        var installedPacks = this.state.installedPacks.map(function(p, idx) {
        return <Pack name={p.Name}
                     id={idx}
                     type={p.Type}
                     version={p.Version}
                     created={p.Created}
                     loadPacks={this.loadPacks}
                     installed={true} />
        }.bind(this));
        var installables = this.state.availablePacks
            .filter(function(p) {
                var installed = false;
                this.state.installedPacks.forEach(function(ip) {
                    if (p.Name == ip.Name && p.Version == ip.Version && p.Created == ip.Created) {
                    installed = true;
                    }
                });
                return !installed;
            }.bind(this))
            .map(function(p, idx) {
                return <Pack name={p.Name}
                             id={idx}
                             type={p.Type}
                             version={p.Version}
                             created={p.Created}
                             file={p.File}
                             loadPacks={this.loadPacks}
                             installed={false}/>
            }.bind(this));
        return <div id="settings-container" onClick={this.hide}>
             <div id="settings-content">
            <div id="settings-installed-packs">
	    <div className="settings-header">Installed Packs</div>
	    <div className="settings-pack-row">
              <div className="settings-pack-row-label">Name</div>
	      <div className="settings-pack-row-label">Type</div>
	      <div className="settings-pack-row-label">Version</div>
	      <div className="settings-pack-row-label">Created</div>
            </div>
	    { installedPacks }
	    </div>
            <div id="settings-available-packs">
	    <div className="settings-header">Available Packs</div>
	    <div className="settings-pack-row">
              <div className="settings-pack-row-label">Name</div>
	      <div className="settings-pack-row-label">Type</div>
	      <div className="settings-pack-row-label">Version</div>
	      <div className="settings-pack-row-label">Created</div>
            </div>
	    { installables }</div>
             </div>
             </div>;
    }
});

var Search = React.createClass({
    search: function(text) {
        this.setState({query: text, selected: 0, results: []});
        this.throttledStreamSearch(text);
    },
    streamSearch: function(text){
        this.setState({query: text, selected: 0, results: []});
        var qs = text.split(" ");
        if (qs.length < 2) {
            return
        }
        this.props.sock.close()
        var pk = qs[0]
        var pt = qs[1]
        var m  = qs[2] || ""
        var req = {
            "Limit": 3,
            "Pack": pk,
            "Path": pt,
            "Member": m
        }

        this.props.sock = socket();
        this.props.sock.onmessage = function(msg) {
            var entry = JSON.parse(msg.data)
            this.setState({results: this.state.results.concat([entry])});
        }.bind(this);
        this.props.sock.onopen = function() {
            this.props.sock.send(JSON.stringify(req));
        }.bind(this);
        this.props.sock.onclose = function() {
            console.log("Finished request [" + text + "].");
        }.bind(this);
    },
    componentWillMount: function() {
        var delay = 80;
        this.throttledStreamSearch = _.debounce(this.streamSearch, delay);
    },
    getInitialState: function(){
        return{
            selected: 0,
            query:'',
            results: []
        }
    },
    selectResult: function(idx) {
        if (idx >= 0 && idx <= 3 && idx < this.state.results.length) {
            this.setState({selected: idx})
        }
    },
    shiftSelection: function(left, right) {
        // if left: try -1
        var sub1 = this.state.selected - 1
        if (left && sub1 >= 0) {
            this.setState({selected: sub1})
            return
        }
        // if right: try +1
        var add1 = this.state.selected + 1
        if (right && add1 < 3 && add1 < this.state.results.length) { // TODO magic number
            this.setState({selected: add1})
            return
        }
        // if none: 0
        if (!left && !right) {
            this.setState({selected: 0})
            return
        }
    },
    componentDidMount: function() {
        key.unbind('shift+right');
        key.unbind('shift+left');

        key('shift+right', function(event) {
            this.shiftSelection(false, true)
            event.cancelBubble = true;
            return false;
        }.bind(this));

        key('shift+left', function(event) {
            this.shiftSelection(true, false)
            event.cancelBubble = true;
            return false;
        }.bind(this));
    },
    showSettings: function () {
        $("#settings-container").css("visibility", "visible");
    },
    render: function(){
        var entries = [];
        for (var i = 0; i < this.state.results.length && i < 3; i++) {
            var entry = this.state.results[i]
            entries.push(<SearchResult
                             entry={entry}
                             index={i}
                             selected={i == this.state.selected}
                             selectResult={this.selectResult} />)
        }

        var params = window.location.search.substring(1);
        var arrParam = params.split("=");

        var doc = "/ui/readme.html"
        if (arrParam.length == 2 && arrParam[0] == "doc") {
            doc = arrParam[1] + decodeURIComponent(window.location.hash).replace(/&gt;/g, ">").replace(/&lt;/g, "<");
        }

        return (<div id="main-container">
                  <Settings />
                  <div id="menu-container"><i className="fa fa-cogs" onClick={this.showSettings}></i></div>
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
    $("#search-field").focus();
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

React.render(<Search sock={socket()} />, document.getElementById("search"));
