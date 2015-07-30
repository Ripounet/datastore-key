(function($) {

	//
	// Button handlers issue requests to services, retrieve JSON data and fill the form fields.
	//
	// Feel free to call /encode and /decode directly for your needs.
	//
	
	$("#ajax-encode").click(function(){
		$("#ajax-encode").button('loading');
		$.ajax({
			url: "/encode",
		    dataType: "json",
			data: $(".form-encode").serialize(),
			success: function(response) {
					var box = $(".form-decode textarea[name=keystring]");
					box.val( response.keystring );
					box.focus();
					$("#ajax-encode").button('reset');
					$(".form-decode").effect("highlight", 1200);
					$(".form-decode fieldset").effect("highlight", 1200);
					window.history.pushState('', '', '/?keystring=' + response.keystring);
			},
			error: function(msg) {
			      alert( "Encoding went wrong : [" + err.responseText + "]" );
				  $("#ajax-encode").button('reset');
			}
		});
	});

	$.fn.ajaxDecode = function(afterSuccess){
		$.ajax({
			url: "/decode",
		    dataType: "json",
			data: $(".form-decode").serialize(),
			success: function(response) {
				$(".form-encode").find("input[type=text], textarea").val("");  // Clear all values
				
				$(".form-encode input[name=kind]").val( response.kind );
				$(".form-encode input[name=intid]").val( response.intID );
				$(".form-encode input[name=stringid]").val( response.stringID );
				$(".form-encode input[name=appid]").val( response.appID );
				$(".form-encode input[name=appid]").change(); // update link
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
				if( afterSuccess )
					afterSuccess();
			},
			error: function(err) {
			      alert( "Key string seems invalid : [" + err.responseText + "]" );
			      $("#ajax-decode").button('reset');
			}
		});
	}
	
	$("#ajax-decode").click(function(){
		$("#ajax-decode").button('loading');
		$.fn.ajaxDecode( function(){
		    $("#ajax-decode").button('reset');
			$(".form-encode").effect("highlight", 1200);
			$(".form-encode fieldset").effect("highlight", 1200);
			window.history.pushState('', '', '/?keystring=' + $(".form-decode textarea[name=keystring]").val());
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
	
	$(".form-encode input[name=appid]").change(function(){
		var appid = $(this).val();
		var url = "javascript:void(0);";
		if( appid ){
			var pos = appid.indexOf("~"); 
			if( pos != -1 )
				appid = appid.substring(pos+1);
			url = "https://" + appid + ".appspot.com";
		}
		$("#appspot-link").attr("href", url);
	});
	$(".form-encode input[name=appid]").change();
	
	$("#btn-more").click(function() {
	    $("#more-content").collapse('toggle');
	});
	
	$.fn.openInDatastoreViewer = function(){
	    var key = $(".form-decode textarea[name=keystring]").val();
	    if( key ){
	    	var kind = $(".form-encode input[name=kind]").val();
	    	var appid = $(".form-encode input[name=appid]").val();
	    	var namespace = $(".form-encode input[name=namespace]").val();
	    	if( appid && kind ){
	    		var url= "https://appengine.google.com/datastore/explorer?submitted=1&app_id=" + appid 
	    			+ "&show_options=yes&viewby=gql&query=SELECT+*+FROM+" + kind 
	    			+ "+WHERE+__key__%3DKEY%28%27"+ key + "%27%29"
	    			+ "&namespace=" + namespace
	    			+ "&options=Run+Query" ;
	    		window.open( url, "datastoreViewer" );
	    	}else{
	    		alert("Please click the Decode button first, to retrieve the App ID.")
	    	}
	    }else{
    		alert("Please provide a datastore key");
    	}
	}
	
	$("#open-in-datastore-viewer").click($.fn.openInDatastoreViewer);
	
	$("#link-for-bookmark").click(function() {
	    var url= window.location.protocol + "//" + window.location.hostname + "/?";
	    [ "kind", "intid", "stringid", "appid", "namespace", "kind2", "intid2", "stringid2", "kind3", "intid3", "stringid3" ].forEach(function(f){
	    	var v = $(".form-encode input[name="+f+"]").val();
	    	if( v )
	    		url += f + "=" + encodeURIComponent(v) + "&";
	    });
	    var keystring =  $(".form-decode textarea[name=keystring]").val();
	    if( keystring )
	    	url += "keystring=" + encodeURIComponent(keystring);
	    window.location = url;
	});
	
	$("#link-engine-this").click(function() {
		window.external.AddSearchProvider( "http://datastore-key.appspot.com/static/xml/opensearch-this.xml" );
	});
	
	$("#link-engine-ds-viewer").click(function() {
		window.external.AddSearchProvider( "http://datastore-key.appspot.com/static/xml/opensearch-jump-to-datastore-viewer.xml" );
	});
	
	
	$("#btn-about").click(function() {
		$("#about-content").collapse('toggle');
	});
})(jQuery);