$(document).ready(function () {
  animateDropdownArrowIcon();
  fillEmptyTableData();
  sortDateRows();
  sortStringRows();
});

var iconUp = `
  <svg width="10" height="16" class="angle-up" viewBox="0 0 320 512" >
    <path fill="currentColor" d="M177 159.7l136 136c9.4 9.4 9.4 24.6 0 33.9l-22.6 22.6c-9.4 9.4-24.6 9.4-33.9 0L160 255.9l-96.4 96.4c-9.4 9.4-24.6 9.4-33.9 0L7 329.7c-9.4-9.4-9.4-24.6 0-33.9l136-136c9.4-9.5 24.6-9.5 34-.1z"></path>
  </svg>`;

var iconRight =
  `<svg width="8" height="16" class="angle-right" viewBox="0 0 256 512">
        <path fill="currentColor" d="M224.3 273l-136 136c-9.4 9.4-24.6 9.4-33.9 0l-22.6-22.6c-9.4-9.4-9.4-24.6 0-33.9l96.4-96.4-96.4-96.4c-9.4-9.4-9.4-24.6 0-33.9L54.3 103c9.4-9.4 24.6-9.4 33.9 0l136 136c9.5 9.4 9.5 24.6.1 34z"></path>
       </svg>`;

var iconDown =
  `<svg width="10" height="16" class="angle-down" viewBox="0 0 320 512">
        <path fill="currentColor" d="M143 352.3L7 216.3c-9.4-9.4-9.4-24.6 0-33.9l22.6-22.6c9.4-9.4 24.6-9.4 33.9 0l96.4 96.4 96.4-96.4c9.4-9.4 24.6-9.4 33.9 0l22.6 22.6c9.4 9.4 9.4 24.6 0 33.9l-136 136c-9.2 9.4-24.4 9.4-33.8 0z"></path>
       </svg>`;

function animateDropdownArrowIcon() {
  $('.report-btn').append(iconRight);

  $('.report-btn').click(function () {
    var expanded = $.parseJSON($(this).attr('aria-expanded'));
    if (!expanded) {
      $(this).children('svg').replaceWith(iconDown);
    } else {
      $(this).children('svg').replaceWith(iconRight);
    }
  });
}

function fillEmptyTableData() {
  $("td").each(function () {
    if (!$(this).text().trim().length) {
      $(this).text("None");
    }
  });
}

function sortDateRows() {
  commonTableSort('.date-th',
    function (a, b) {
      var aDate = $(a).children(".date-td").text();
      var bDate = $(b).children(".date-td").text();

      return new Date(aDate) - new Date(bDate);
    },
    function (a, b) {
      var aDate = $(a).children(".date-td").text();
      var bDate = $(b).children(".date-td").text();

      return new Date(bDate) - new Date(aDate);
    });
}

function sortStringRows() {
  commonTableSort('.string-th',
    function (a, b) {
      var aStr = $(a).children('.string-td').text();
      var bStr = $(b).children('.string-td').text();

      return aStr.localeCompare(bStr);
    },
    function (a, b) {
      var aStr = $(a).children('.string-td').text();
      var bStr = $(b).children('.string-td').text();

      return bStr.localeCompare(aStr);
    });
}

function commonTableSort(th, ascSortFunc, descSortFunc) {
  // set initial arrow icon
  $(th).append(iconDown);
  $(th).children('svg').css("margin-left", "10px");

  $(th).click(function () {
    var parentTable = $(this).parents()[2];
    var parentTableHead = $(parentTable).children('thead');
    var tableHeaders = $($(parentTableHead).children('tr')).children('th');

    var parentTableBody = $(parentTable).children('tbody');
    var tableRows = $(parentTableBody).children('tr').toArray();

    var sorted = $.parseJSON($(this).attr('sorted'));

    if (!sorted) {
      tableRows.sort(ascSortFunc);
      // set right icon on click
      $(this).children('svg').replaceWith(iconDown);
      $(this).children('svg').css("margin-left", "10px");
      $(this).attr('sorted', true);
    } else {
      tableRows.sort(descSortFunc);
      // set right icon on click
      $(this).children('svg').replaceWith(iconUp);
      $(this).children('svg').css("margin-left", "10px");
      $(this).attr('sorted', false);
    }

    // reset all other headers
    $(tableHeaders).not(this).attr('sorted', false);

    // sort indexes that count table rows
    organizeIndexes(tableRows);
    // append sorted tablerows to table
    $(parentTableBody).append(tableRows);
  });
}

function organizeIndexes(tableRows) {
  var i = 1;
  tableRows.forEach(x => {
    $(x).children("th").text(i.toString());
    i++;
  });
}
