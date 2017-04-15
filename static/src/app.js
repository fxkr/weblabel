require('file-loader?name=[name].[ext]!./index.html');

import Bootstrap from 'bootstrap/dist/css/bootstrap.css';
require('./app.styl')

var target = "http://localhost:8081";

function onPrintButtonClick(e) {
  e.preventDefault();

  $.post({
    url: target + "/api/v1/printer/print",
    data: JSON.stringify({text: $("#label-text").val()}),
    dataType: "json",
  }).done(function(data) {
    console.log("done", data);
  });
}

$().ready(() => {
  $("#print-button").click(onPrintButtonClick);
});
