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
        this.state.results.forEach(function(r) {
            var target = "/" + r["Target"];
            entries.push(<div onClick={this.loadDoc.bind(this, target)}><a href={target}>{r["Entity"]} {r["Function"]}</a> {r["Signature"]}</div>)
        }.bind(this));

        return (<div>
                  <SearchField query={this.state.query} search={this.search}/>
                  <div>{entries}</div>
                  <div>
                    <iframe id="ifrm" src="/" style={{width:'1000px',height:'400px'}} />
                  </div>
                </div>
        );
    }
});

React.render(<Search />, document.body);
