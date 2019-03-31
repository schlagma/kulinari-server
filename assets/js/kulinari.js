if (document.getElementById('btnSidebar') != null) {
    document.getElementById('btnSidebar').addEventListener('click', function(){
        var sidebar = document.getElementById('sidebarContent');
        var overlay = document.getElementById('overlay');
        if(sidebar.style.visibility === 'visible') {
            sidebar.style.transform = 'translateX(-20rem)';
            sidebar.style.visibility = 'hidden';
            overlay.style.visibility = 'hidden';
            overlay.style.backgroundColor = 'transparent';
        } else {
            sidebar.style.transform = 'translateX(0)';
            sidebar.style.visibility = 'visible';
            overlay.style.visibility = 'visible';
            overlay.style.backgroundColor = 'rgba(0,0,0,0.8)';
        }
    });
}

if (document.getElementById('overlay') != null) {
    document.getElementById('overlay').addEventListener('click', function(){
        var sidebar = document.getElementById('sidebarContent');
        var overlay = document.getElementById('overlay');
        if(sidebar.style.visibility === 'visible') {
            sidebar.style.transform = 'translateX(-20rem)';
            sidebar.style.visibility = 'hidden';
            overlay.style.visibility = 'hidden';
            overlay.style.backgroundColor = 'transparent';
        } else {
            sidebar.style.transform = 'translateX(0)';
            sidebar.style.visibility = 'visible';
            overlay.style.visibility = 'visible';
            overlay.style.backgroundColor = 'rgba(0,0,0,0.8)';
        }
    });
}

if (document.getElementById('btnMenu') != null) {
    document.getElementById('btnMenu').addEventListener('click', function(){
        var menu = document.getElementById('menu');
        if(menu.style.display === 'inline') {
            menu.style.display = 'none';
        } else {
            menu.style.display = 'inline';
        }
    });
    window.addEventListener('mouseup', function(event){
        var menu = document.getElementById('menu');
        var btnMenu = document.getElementById('btnMenu');
        if(event.target != menu || event.target == btnMenu) {
            menu.style.display = 'none';
        }
    });
}

if (document.getElementById('btnCatBack') != null) {
    document.getElementById('btnCatBack').addEventListener('click', function(){ window.history.back(); });
}

if (document.getElementById('btnPrint') != null) {
    document.getElementById('btnPrint').addEventListener('click', function(){ window.print(); });
}

if (document.getElementById('ingCalc') != null) {
    document.getElementById('ingCalc').addEventListener('click', function(){
        var x = document.getElementById('ingVal').value;
        var y = window.location.href;
        var y = y.split('?');
        var y = y[0];
        if (x > 0) {
            window.location = y + '?for=' + x;
        }
    });
}

if (document.getElementById('recFilter') != null) {
    document.getElementById('recFilter').addEventListener('keyup', function(){
        var input, filter, ul, li, a, i;
        input = document.getElementById('recFilter');
        filter = input.value.toUpperCase();
        ul = document.getElementById('recUl');
        li = ul.getElementsByTagName('li');

        for (i = 0; i < li.length; i++) {
            a = li[i].getElementsByTagName("a")[0];
            if (a.innerHTML.toUpperCase().indexOf(filter) > -1) {
                li[i].style.display = "";
            } else {
                li[i].style.display = "none";
            }
        }

        a = li[0].getElementByTagName('a');
        a.focus();
    });
}

if (document.getElementById('catFilter') != null) {
    document.getElementById('catFilter').addEventListener('keyup', function(){
        var input, filter, ul, li, a, i;
        input = document.getElementById('catFilter');
        filter = input.value.toUpperCase();
        ul = document.getElementById('cookbooks');
        li = ul.getElementsByTagName('a');

        for (i = 0; i < li.length; i++) {
            a = li[i].getElementsByTagName('h1')[0];
            if (a.innerHTML.toUpperCase().indexOf(filter) > -1) {
                li[i].style.display = "";
            } else {
                li[i].style.display = "none";
            }
        }
    });
}

if (document.getElementById('pRecFilter') != null) {
    document.getElementById('pRecFilter').addEventListener('keyup', function(){
        var input, filter, ul, li, a, i;
        input = document.getElementById('pRecFilter');
        filter = input.value.toUpperCase();
        ul = document.getElementById('pRecUl');
        li = ul.getElementsByTagName('div');

        for (i = 0; i < li.length; i++) {
            a = li[i].getElementsByTagName('span')[0];
            if (a.innerHTML.toUpperCase().indexOf(filter) > -1) {
                li[i].style.display = "";
            } else {
                li[i].style.display = "none";
            }
        }
    });
}

if (document.getElementById('btnDelete') != null) {
    document.getElementById('btnDelete').addEventListener('click', function(){
        var modal = document.getElementById('modalDelete');
        if(modal.style.display === 'inline') {
            modal.style.display = 'none';
        } else {
            modal.style.display = 'inline';
        }
    });
    document.getElementById('modalDeleteClose').addEventListener('click', function(){
        var modal = document.getElementById('modalDelete');
        if(modal.style.display === 'none') {
            modal.style.display = 'inline';
        } else {
            modal.style.display = 'none';
        }
    });
    document.getElementById('modalDeleteCancel').addEventListener('click', function(){
        var modal = document.getElementById('modalDelete');
        if(modal.style.display === 'none') {
            modal.style.display = 'inline';
        } else {
            modal.style.display = 'none';
        }
    });
}

if (document.getElementById('btnTimer') != null) {
    document.getElementById('btnTimer').addEventListener('click', function(){
        var modal = document.getElementById('modalTimer');
        if(modal.style.display === 'inline') {
            modal.style.display = 'none';
        } else {
            modal.style.display = 'inline';
        }
    });
    document.getElementById('modalTimerClose').addEventListener('click', function(){
        var modal = document.getElementById('modalTimer');
        if(modal.style.display === 'none') {
            modal.style.display = 'inline';
        } else {
            modal.style.display = 'none';
        }
    });
}

if (document.getElementById('btnStartStop') != null) {
    document.getElementById('btnStartStop').addEventListener('click', function(){
        
        var hour = document.getElementById("hourInput").value;
        var minute = document.getElementById("minuteInput").value;
        var second = document.getElementById("secondInput").value;

        var time = (hour * 3600) + (minute * 60) + second;

        if (hour < 10) {
            document.getElementById('hour').innerHTML = "0" + hour;
        } else {
            document.getElementById('hour').innerHTML = hour;
        }

        if (minute < 10) {
            document.getElementById('minute').innerHTML = "0" + minute;
        } else {
            document.getElementById('minute').innerHTML = minute;
        }

        if (second < 10) {
            document.getElementById('second').innerHTML = "0" + second;
        } else {
            document.getElementById('second').innerHTML = second;
        }

        const update = () => {
            while (time != 0) {
                time--;
                hour = parseInt(time / 3600);
                minute = parseInt((time / 60) - (hour * 60));
                second = parseInt(time - (hour * 3600) - (minute * 60));

                if (hour < 10) {
                    document.getElementById('hour').innerHTML = "0" + hour;
                } else {
                    document.getElementById('hour').innerHTML = hour;
                }
        
                if (minute < 10) {
                    document.getElementById('minute').innerHTML = "0" + minute;
                } else {
                    document.getElementById('minute').innerHTML = minute;
                }
        
                if (second < 10) {
                    document.getElementById('second').innerHTML = "0" + second;
                } else {
                    document.getElementById('second').innerHTML = second;
                }
            }
        };
        setInterval(update, 1000);
    });
}

if (document.getElementById('toTop') != null) {
    document.getElementById('toTop').addEventListener('click', function(){
        scrollArea = document.getElementById('scrollArea');
        scrollArea.scrollTo(0, 0);
    });
}
