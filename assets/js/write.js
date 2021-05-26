$(function(){
    $("#content").bind("input change", function(){
        $.post("/getHtml", {md: $("#content").val()}, function(response){
            $("#md_html").html(response.html)
        })
    })
})