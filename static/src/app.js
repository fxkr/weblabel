import 'file-loader?name=[name].[ext]!./index.html';

import 'bootstrap/dist/css/bootstrap.css';
import 'github-fork-ribbon-css/gh-fork-ribbon.css';
import './app.styl';


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
