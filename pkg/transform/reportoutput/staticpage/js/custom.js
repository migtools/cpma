$(document).ready(function () {
  animateArrowIcon();
  fillEmptyTableData();
});

function animateArrowIcon() {
  var iconRight =
    `<svg width="8" height="16" class="angle-right" viewBox="0 0 256 512">
        <path fill="currentColor" d="M224.3 273l-136 136c-9.4 9.4-24.6 9.4-33.9 0l-22.6-22.6c-9.4-9.4-9.4-24.6 0-33.9l96.4-96.4-96.4-96.4c-9.4-9.4-9.4-24.6 0-33.9L54.3 103c9.4-9.4 24.6-9.4 33.9 0l136 136c9.5 9.4 9.5 24.6.1 34z"></path>
       </svg>`;

  var iconDown =
    `<svg width="10" height="16" class="angle-down" viewBox="0 0 320 512">
        <path fill="currentColor" d="M143 352.3L7 216.3c-9.4-9.4-9.4-24.6 0-33.9l22.6-22.6c9.4-9.4 24.6-9.4 33.9 0l96.4 96.4 96.4-96.4c9.4-9.4 24.6-9.4 33.9 0l22.6 22.6c9.4 9.4 9.4 24.6 0 33.9l-136 136c-9.2 9.4-24.4 9.4-33.8 0z"></path>
       </svg>`;

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
