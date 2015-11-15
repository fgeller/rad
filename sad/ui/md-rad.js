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
							<input className="mdl-textfield__input" type="text" id="search-field-input" />
							<label className="mdl-textfield__label" htmlFor="search-field-input"><span id="search-label-pack">pack</span> path member</label>
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
