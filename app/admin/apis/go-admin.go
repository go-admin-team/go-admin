package apis

import (
	"github.com/gin-gonic/gin"
)

const INDEX = `
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>GO-ADMIN欢迎您</title>
<style>
body{
  margin:0; 
  padding:0; 
  overflow-y:hidden
}
</style>
<script src="http://libs.baidu.com/jquery/1.9.0/jquery.js"></script>
<script type="text/javascript"> 
window.onerror=function(){return true;} 
$(function(){ 
  headerH = 0;  
  var h=$(window).height();
  $("#iframe").height((h-headerH)+"px"); 
});
</script>
</head>
<body>
<iframe id="iframe" frameborder="0" src="https://doc.go-admin.dev" style="width:100%;"></iframe>
</body>
</html>
`

func GoAdmin(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(200, INDEX)
}
