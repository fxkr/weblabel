require('file-loader?name=[name].[ext]!./index.html');

import Bootstrap from 'bootstrap/dist/css/bootstrap.css';
require('./app.styl')

function onPrintButtonClick(e) {
  e.preventDefault();

  $.post({
    url: "/api/v1/printer/print",
    data: JSON.stringify({text: $("#label-text").val()}),
    dataType: "json",
  }).done(function(data) {
    console.log("done", data);
  });
}

$().ready(() => {
  $("#print-button").click(onPrintButtonClick);
});
