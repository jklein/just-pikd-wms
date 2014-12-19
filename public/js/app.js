// This file bootstraps the entire application.

var WmsApp = require('./components/WmsApp.react');
var React = require('react');
window.React = React; // export for http://fb.me/react-devtools

React.render(
  <WmsApp />,
  document.getElementById('react')
);




