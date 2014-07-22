$(function() {

	$("#ajax-encode").click(function(){
		$.getJSON("/encode",
				$(".form-encode").serialize(),
				function(response, status) {
			alert(status);
					var box = $(".form-decode textarea[name=keystring]");
					box.val( response.keystring );
					box.focus();
					$(".form-decode").effect("highlight", 1000);
				});
	});

	$("#ajax-decode").click(function(){
		$.getJSON("/decode",
				$(".form-decode").serialize(),
				function(response, status) {
			alert(status);
					$(".form-encode").find("input[type=text], textarea").val("");  // Clear all values
			
					$(".form-encode input[name=kind]").val( response.kind );
					$(".form-encode input[name=intid]").val( response.intID );
					$(".form-encode input[name=stringid]").val( response.stringID );
					$(".form-encode input[name=appid]").val( response.appID );
					$(".form-encode input[name=namespace]").val( response.namespace );
					
					if( response.parent ){
						setSectionVisibility( $("#set-parent"), $(".key-parent"), true, "Remove parent" );
						$(".form-encode input[name=kind2]").val( response.parent.kind );
						$(".form-encode input[name=intid2]").val( response.parent.intID );
						$(".form-encode input[name=stringid2]").val( response.parent.stringID );
						
						if( response.parent.parent ){
							setSectionVisibility( $("#set-grand-parent"), $(".key-grand-parent"), true, "Remove grandparent" );
							$(".form-encode input[name=kind3]").val( response.parent.parent.kind );
							$(".form-encode input[name=intid3]").val( response.parent.parent.intID );
							$(".form-encode input[name=stringid3]").val( response.parent.parent.stringID );
						}else{
							setSectionVisibility( $("#set-grand-parent"), $(".key-grand-parent"), false, "Set grandparent" );
						}
					}else{
						setSectionVisibility( $("#set-parent"), $(".key-parent"), false, "Set parent" );
						setSectionVisibility( $("#set-grand-parent"), $(".key-grand-parent"), false, "Set grandparent" );
					}
					$(".form-encode").effect("highlight", 1000);
				});
	});

	//
	// Toggle Set/Remove the direct parent of the main key
	//
	$("#set-parent").click(function(){
		if( $(".key-parent").is(":visible") ){
			$(".form-encode .key-parent").find("input[type=text], textarea").val("");
			setSectionVisibility( $(this), $(".key-parent"), false, "Set parent" );
			setSectionVisibility( $(this), $(".key-grand-parent"), false, "Set grandparent" );
		}else{
			setSectionVisibility( $(this), $(".key-parent"), true, "Remove parent" );
		}
	});

	//
	// Toggle Set/Remove the parent of the parent of the main key
	//
	$("#set-grand-parent").click(function(){
		if( $(".key-grand-parent").is(":visible") ){
			$(".form-encode .key-grand-parent").find("input[type=text], textarea").val("");
			setSectionVisibility( $(this), $(".key-grand-parent"), false, "Set grandparent" );
		}else{
			setSectionVisibility( $(this), $(".key-grand-parent"), true, "Remove grandparent" );
		}
	});
	
	function setSectionVisibility(button, section, newVisibility, newButtonText){
		if( newVisibility ){
			section.removeClass("hidden");
			section.focus();
		}else{
			section.addClass("hidden");
		}
		button.html(newButtonText);
	}


	//
	// There is no support for great-grand-parents and further ancestors, but contact me if you feel you need that.
	//
});