var SearchField = React.createClass({
    search: function() {
        var query = this.refs.search.getDOMNode().value;
        this.props.search(query)
    },
    render: function() {
        return <input
                 type="text"
                 ref="search"
                 placeholder="Search me..."
                 value={this.props.query}
                 onChange={this.search} />
    }
});

var Search = React.createClass({
    search: function(text){
        this.setState({query: text});
        var qs = text.split(" ");
        if (qs.length > 1) {
            var q = "/s?p=" + qs[0] + "&e=" + qs[1]
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
    loadDoc: function(target) {
        console.log("target", target);
        document.getElementById("ifrm").src = target;
        this.setState({loadedTarget: target});
    },
    render: function(){
        var entries = [];
        for (var i = 0; i < this.state.results.length && i < 5; i++) {
            var entry = this.state.results[i]
            var target = "/packs/" + entry["Target"];
            entries.push(<div onClick={this.loadDoc.bind(this, target)}>{entry["Entity"]} {entry["Function"]} {entry["Signature"]}</div>)
        }

        return (<div>
                  <SearchField query={this.state.query} search={this.search}/>
                  <div>{entries}</div>
                  <div>
                  <iframe id="ifrm" src="" />
                    </div>
                  </div>
        );
    }
});

React.render(<Search />, document.body);
