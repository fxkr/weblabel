import 'file-loader?name=[name].[ext]!./index.html';

import 'bootstrap/dist/css/bootstrap.css';
import 'github-fork-ribbon-css/gh-fork-ribbon.css';
import 'noty/lib/noty.css';
import './app.styl';

import Noty from 'noty';

function onPrintButtonClick(e) {
  e.preventDefault();

  var note = new Noty({
    text: 'Printing...',
    type: 'info',
    layout: 'bottomRight',
  }).show();

  $.post({
    url: "/api/v1/printer/print",
    data: JSON.stringify({text: $("#label-text").val()}),
    dataType: "json",
  }).done((data) => {
    note.setText("Label printed successfully.");
    note.setType('success');
    note.setTimeout(1000);
  }).fail((data) => {
    note.setText("Failed to print label.");
    note.setType('error');
    note.setTimeout(3000);
  });
}

$().ready(() => {
  $("#print-button").click(onPrintButtonClick);
});
