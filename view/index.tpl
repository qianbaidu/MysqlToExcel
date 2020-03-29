<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta name="csrf-param" content="_csrf">
    <meta name="csrf-token" content="Y3NzZVRaa1rmSDyZaYAvHJoMKesdSMOfWTVqGxZsme60CDWLBezq.w==">
    <title>MysqlToExcel</title>
    <link href="static/css/H-ui.css" rel="stylesheet">
    <style>
        form > .form-group {
            display: block;
        }

        h1 {
            display: none;
        }
    </style>
</head>
<body>

<div class="wrap">


    <div class="container">
        <article class="page-container">
            <form class="form form-horizontal" id="form-article-add">

                <div class="row cl">
                    <label class="form-label col-xs-4 col-sm-2"> save file name：</label>


                    <div class="formControls col-xs-8 col-sm-8">
                        <input id="name" type="text" class="input-text" value="" placeholder="" name="">
                    </div>
                </div>




                <div class="row cl">
                    <label class="form-label col-xs-4 col-sm-2"><span class="c-red">*</span> db：</label>
                    <div class="formControls col-xs-8 col-sm-8">
                        <input id="db" type="text" class="input-text" value="Excel" placeholder="" name="">
                    </div>
                </div>

                <div class="row cl">
                    <label class="form-label col-xs-4 col-sm-2">sql：</label>
                    <div class="formControls col-xs-8 col-sm-9">
                        <textarea name="sql" id="sql-val" cols="" rows="" class="textarea" placeholder="sql"
                                  datatype="*10-100" dragonfly="true" nullmsg="sql不能为空！"></textarea>

                    </div>
                </div>


                <div class="row cl">
                    <div class="col-xs-8 col-sm-9 col-xs-offset-4 col-sm-offset-2">
                        <button id="putout" class="btn btn-secondary radius" type="button">导出</button>
                        <button class="btn btn-default radius" type="reset">&nbsp;&nbsp;取消&nbsp;&nbsp;</button>
                    </div>
                </div>

                <br/>
                <br/>
                <hr/>

                <div class="row cl" id="saveName-content" style="display: none">
                    <label class="form-label col-xs-4 col-sm-2">输出文件：</label>
                    <div class="formControls col-xs-8 col-sm-9">
                        <input id="savedName" type="text" class="input-text" value="" style="border: none;">
                    </div>
                </div>
            </form>
        </article>

    </div>
</div>


<script src="static/js/jquery.js"></script>
<script type="text/javascript">jQuery(document).ready(function () {

        var putOutUrl = "http://localhost:9010/MysqlToExcel";
        $(function () {

            window.alertSuccess = {};
            window.alertError = {};
            window.appendStr = {};

            function success(str) {
                $('.data-alert').animate({top: "0px", display: 'block'}, 1);
                $('.date-success-alert').html(str).show().animate({top: "50px"}, 500);
                var timer = setTimeout(function () {
                    $('.date-success-alert').animate({top: "-100px", display: 'none'}, 800);
                    $('.data-alert').animate({top: "-100px", display: 'none'}, 1000);
                    clearTimeout(timer);
                }, 3000);

            }

            function error(str) {
                $('.data-alert').animate({top: "0px", display: 'block'}, 1);
                $('.date-error-alert').html(str).show().animate({top: "50px"}, 500);
                var timer = setTimeout(function () {
                    $('.date-error-alert').animate({top: "-100px", display: 'none'}, 1000);
                    $('.data-alert').animate({top: "-100px", display: 'none'}, 1000);
                    clearTimeout(timer);
                }, 3000);
            }

            function appendStr() {
                str = '<div class="" style="position: fixed;top:0px;left:0px;padding:0px;width:100%;height:0px;z-index: 999999;">' +
                    '<div class="data-alert" style="margin: 0 auto;width: 50%;min-width: 200px;padding: 0px ;position: relative;">' +
                    '<div class="date-success-alert" style="padding: 0px 50px;display:none;box-shadow: 5px 5px 5px rgb(221, 221, 221);font-size: 16px;background: #fff;text-align: center;height: 38px;line-height: 38px;position: relative;top: -100px;border-radius: 5px; color: #fff;background: rgba(0,192,254,0.8);">操作成功' +
                    '</div>' +
                    '<div class="date-error-alert" style="padding: 0px 50px;display:none;background:box-shadow: 5px 5px 5px rgb(221, 221, 221); font-size: 16px;#fff;text-align: center;height: 38px;line-height: 38px;position: relative;top: -150px;border-radius: 5px; color: yellow;background: rgba(255,0,0,0.8);">操作失败</div>' +
                    '</div> </div>';
                $('body').append(str);
            }

            window.alertSuccess = success;
            window.alertError = error;
            window.appendStr = appendStr;

            window.appendStr();


            function ajaxPost(url, postData, callBack) {
                var ajax = $.ajax({
                    type: "POST",
                    url: url,
                    data: postData,
                    timeout: 15000,
                    success: callBack,
                    error: function (XMLHttpRequest, textStatus, errorThrown) {
                        if (textStatus == 'error' || errorThrown == 'Internal Server Error') {
                            window.alertError('服务器异常');
                        }
                    },
                    complete: function (XMLHttpRequest, status) {
                        if (status == 'timeout') {
                            ajax.abort();
                            window.alertError('请求超时，请稍后重试');
                        }
                    }
                });

            }

            //导出
            $("#putout").click(function () {
//            var newTab = window.open('about:blank');
                $('#saveName-content').hide();
                var sql = $("#sql-val").val();
                var db = $("#db").val();
                var name = $("#name").val();
                if (sql == '') {
                    window.alertError("sql不能为空");
                    return false;
                }
                var postData = "sql=" + sql + "&db=" + db + "&type=1&name=" + name;

                ajaxPost(putOutUrl, postData, function (data) {
                    obj = JSON.parse(data);
                    console.log(obj);
                    if (obj.code == 10000) {
                        window.alertSuccess("导出完成");
                        $('#saveName-content').show();
                        $('#savedName').val(obj.name);
                    } else {
                        window.alertError("导出失败");
                    }
                })
            })

        })
    });
</script>
</body>
</html>
