var data = {}


function highlightMatch(source, match){
  var titleLower = source.toLowerCase();
  var start_pos = titleLower.search(match.toLowerCase());
  var end_pos = start_pos+match.length;
  result = source.substring(0, start_pos)+'<b>'+ source.substring(start_pos,end_pos)+'</b>'+source.substring(end_pos,source.length);
  return result;
}

$('document').ready(function(){


  /*
  * Keep scrolling position of navigation
  */

  if(localStorage.getItem('scroll_position') != null){
    $( '.sidebar' ).scrollTop(localStorage.getItem('scroll_position'));
  }

  $( '.sidebar' ).scroll(function() {
      localStorage.setItem('scroll_position',$(this).scrollTop());
  });


  /*
  * Sidebar Search
  */

  var sidebar = $('#sidebar-list').html();
  $.getJSON("search.json", function(json) {
      data = json;
  });
  $( "#searchbox" ).keyup(function(){
    if($(this).val() == ''){
      $('#sidebar-list').html(sidebar);
    }else{
      var input = $(this).val();
      var out = '<ul>';
      var foundAnything = false;
      var foundIn = [];
      $.each(data,function(key,entry) {
        var contains = false;
        // check if title contains search
        if(entry.title.toLowerCase().indexOf(input.toLowerCase()) !== -1){
          contains = true;
          foundIn.push('title');
        }
        //check if tags contain search
        for(var i = 0; i < entry.tags.length; i++){
          if(entry.tags[i].toLowerCase().indexOf(input.toLowerCase()) !== -1){
            contains = true;
            foundIn.push('tags');
          }
        }
        //check if aliases contain search
        for(var i = 0; i < entry.aliases.length; i++){
          if(entry.aliases[i].toLowerCase().indexOf(input.toLowerCase()) !== -1){
            contains = true;
            foundIn.push('entry');
          }
        }

        // render output
        if(contains == true){
          foundAnything = true;
          out += '<li>';
          var display = this.title;

          if(foundIn.indexOf('title') > -1){
            // set title bold if keyword is in it
            display = highlightMatch(this.title,input);
          }

          out += '<a href="'+key+'">'+display+'</a>';
          if(foundIn.indexOf('entry') > -1){
            for(var i = 0; i < entry.aliases.length; i++){
              if(entry.aliases[i].toLowerCase().indexOf(input.toLowerCase()) !== -1){
                out += '<div class="search-alias">Alias: ';
                out += highlightMatch(entry.aliases[i],input);
                out += '</div>';
              }
            }
          }

          if(foundIn.indexOf('tags') > -1){
            for(var i = 0; i < entry.tags.length; i++){
              if(entry.tags[i].toLowerCase().indexOf(input.toLowerCase()) !== -1){
                out += '<div class="search-alias">Stichwort: ';
                out += highlightMatch(entry.tags[i],input);
                out += '</div>';
              }
            }
          }


          out += '</li>';
        }


      });
      if(foundAnything == false){
        out = '<i>No search results found.</i>';
      }else{
        out += '</ul>';
      }
      $('#sidebar-list').html(out);
    }
  });



  /*
  * Mobile Design
  */
  if (matchMedia) {
    var mq = window.matchMedia("(max-width: 768px)");
    mq.addListener(WidthChange);
    WidthChange(mq);
  }

  // media query change
  function WidthChange(mq) {
    if (mq.matches) {
      $('.sidebar').addClass('hidden');
      $('.nav-toggle').removeClass('hidden');

    } else {
      $('.sidebar').removeClass('hidden');
      $('.content').removeClass('hidden');
      $('.nav-toggle').addClass('hidden');
    }
  }

  $('.nav-toggle').click(function(){
    $('.sidebar').toggleClass('hidden');
    $('.content').toggleClass('hidden');
  });
});
