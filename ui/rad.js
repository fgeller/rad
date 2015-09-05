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
                 placeholder="Search me..."
                 value={this.props.query}
                 onChange={this.search} />
    }
});

var SearchResult = React.createClass({
    open: function() {
        var href = "/packs/" + this.props.entry["Target"];
        console.log("opening target", href);
        document.getElementById("ifrm").src = href;
    },
    render: function() {
        var entName = this.props.entry["Entity"];
        var funName = this.props.entry["Function"] || "\u00a0"; // &nbsp;
        return <div className="search-result" onClick={this.open.bind(this)}>
                 <div className="entity-name">{entName}</div>
                 <div className="function-name">{funName}</div>
               </div>
    }
});

var Search = React.createClass({
    search: function(text){
        this.setState({query: text});
        var qs = text.split(" ");
        if (qs.length > 1) {
            var q = "/s?p=" + qs[0] + "&e=" + qs[1]
            if (qs.length > 2) {
                q += "&f=" + qs[2]
            }
            $.get(q, function(result) {
                this.setState({results: result});
            }.bind(this));
        }
    },
    getInitialState: function(){
        return{
            query:'',
            results: []
        }
    },
    render: function(){
        var entries = [];
        for (var i = 0; i < this.state.results.length && i < 5; i++) {
            var entry = this.state.results[i]
            entries.push(<SearchResult entry={entry} />)
        }

        return (<div>
                <div id="search-field-container">
                  <SearchField query={this.state.query} search={this.search}/>
                </div>
                <div id="search-result-container">{entries}</div>
                <div>
                <iframe id="ifrm" src="" />
                  </div>
                </div>
        );
    }
});

React.render(<Search />, document.body);
