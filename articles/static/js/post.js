$(document).ready(function(){
  $(".editable").hallo({
    plugins: {
      'halloformat': {},
      'halloheadings': {},
      'halloreundo': {}
    },
    toolbar: 'halloToolbarFixed'
  });
  $(".blog-main").on('hallomodified', function(c, d) {
    var s = d.content;
    s = s.replace(/<\/?i>/g,"_").replace(/<\/?b>/g,"*")
          .replace(/<div>/g,"\n").replace(/<\/div>/g,"")
          .replace(/<p>/g,"\n").replace(/<\/p>/g,"")
          .replace(/<span[^>]*>/g,"\n#").replace(/<\/span>/g,"")
          .replace(/<h\d[^>]*>/g,"\n#").replace(/<\/h\d>/g,"")
          .replace(/<br>/g,"\n");
    $("#mdsource").val(s);
  })
  $(".blog-title").on('hallomodified', function(c, d) {
    var s = d.content;
    s = s.replace(/<\/?i>/g,"_").replace(/<\/?b>/g,"*")
          .replace(/<div>/g,"\n").replace(/<\/div>/g,"")
          .replace(/<p>/g,"\n").replace(/<\/p>/g,"")
          .replace(/<span[^>]*>/g,"\n#").replace(/<\/span>/g,"")
          .replace(/<h\d[^>]*>/g,"\n#").replace(/<\/h\d>/g,"")
          .replace(/<br>/g,"\n").replace(/\n/g," ");
    $("#mdtitle").val(s);
  })
});
