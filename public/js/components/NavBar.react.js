var React = require('react');
var NavBar = React.createClass ({
    getInitialState: function() {
      return {text: ''};
    },

    render: function() {
        return (
            <nav className="navbar navbar-inverse" role="navigation">
                <div className="container-fluid">
                  <div className="navbar-header">
                    <button type="button" className="navbar-toggle collapsed" data-toggle="collapse" data-target="#bs-example-navbar-collapse-1"></button>
                    <a className="navbar-brand" href="#">Ella{"'"}s Market WMS</a>

                  </div>
                  <button type="button" className="active btn btn-primary navbar-btn">Purchasing</button>
                  <button type="button" className="btn btn-primary navbar-btn">Receiving</button>
                  <button type="button" className="btn btn-primary navbar-btn">Stocking</button>
                  <button type="button" className="btn btn-primary navbar-btn">Picking</button>
                  <button type="button" className="btn btn-primary navbar-btn">Pickup</button>
                  <button type="button" className="btn btn-danger navbar-btn pull-right">Reset Data</button>
                </div>
            </nav>
        );
    }
});

module.exports = NavBar;