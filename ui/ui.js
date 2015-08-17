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
    getInitialState:function(){
        return{
            query:'',
            results: []
        }
    },
    render:function(){
        var entries = [];
        this.state.results.forEach(function(r) {
            entries.push(<div><code>{r}</code></div>)
        });

        return (<div>
                  <SearchField query={this.state.query} search={this.search}/>
                  <div>{entries}</div>
                </div>
        );
    }
});

React.render(<Search />, document.body);
