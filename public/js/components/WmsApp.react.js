var NavBar = require('./Navbar.react');

var React = require('react');

var WmsApp = React.createClass({

  render: function() {
    return (
      <div className="wmsapp">
        <NavBar />
        <h1>Hello!</h1>
      </div>
    );
  }

});

module.exports = WmsApp;
