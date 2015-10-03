var SearchField = React.createClass({
    search: function() {
        var query = this.refs.search.getDOMNode().value;
        this.props.search(query)
    },
    render: function() {
        return <input
                 id="search-field"
                 type="text"
                 ref="search"
                 placeholder="Search here..."
                 value={this.props.query}
                 onChange={this.search} />
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
        var namespace = this.props.entry["Namespace"].join(".") || "\00a0";
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

var Search = React.createClass({
    search: function(text){
        this.setState({query: text, selected: 0});
        var qs = text.split(" ").map(encodeURIComponent);
        if (qs.length > 1) {
            var q = "/s?p=" + qs[0] + "&e=" + qs[1]
            if (qs.length > 2) {
                q += "&m=" + qs[2]
            }
            $.get(q, function(result) {
                this.setState({results: result});
            }.bind(this));
        }
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

        var doc = ""
        if (arrParam.length == 2 && arrParam[0] == "doc") {
            doc = arrParam[1] + decodeURIComponent(window.location.hash).replace(/&gt;/g, ">").replace(/&lt;/g, "<");
        }

        return (<div id="main-container">
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

React.render(<Search />, document.getElementById("search"));
