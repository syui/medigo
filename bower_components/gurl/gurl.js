jQuery(function ($) {
	$.ajax({
		async: true,
		cache: true,
		success: function(data) {
			var pel = document.getElementById("param") ;
			var url = $.url(window.location.href).param();
			console.log(url.code);
			$('#param').text(url.code);     
		}
	});
});
