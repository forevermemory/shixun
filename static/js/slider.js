/*	window.onload = function() {
	var oActionBlock = document.getElementById('action-block');
	var oActionBar = document.getElementById('action-bar');
	var oScrollBar = document.getElementById('scroll-bar');
	var oShowAmount = document.getElementById('showAmount').getElementsByTagName('input')[0];
	var length = 550;

	clickSlide(oActionBlock, oActionBar, oScrollBar, 300, length, oShowAmount);
	drag(oActionBlock, oActionBar, 300, length, oShowAmount);
	addScale(60, 300, length, oScrollBar);
	inputBlur(oActionBlock, oActionBar, length, oShowAmount);
}		*/

function SlideBar(data){
	var _this = this;
	var oActionBlock = document.getElementById(data.actionBlock);
	var oActionBar = document.getElementById(data.actionBar);
	var oSlideBar = document.getElementById(data.slideBar);
	var barLength = data.barLength;
	var interval = data.interval;
	var maxNumber = data.maxNumber;
	var minNumber = data.minNumber;
	var oShowArea = null;
	if(data.showArea){
		oShowArea = document.getElementById(data.showArea);	
	}

	if(oShowArea){
		if(oShowArea.value) {
			oActionBar.style.width = oShowArea.value * 4 + 'px';
			oActionBlock.style.left = oActionBar.offsetWidth - (oActionBlock.offsetWidth / 2) + 'px';
		}
		_this.addScale(oSlideBar, interval,minNumber, maxNumber, barLength);
		_this.clickSlide(oActionBlock, oActionBar, oSlideBar,minNumber, maxNumber, barLength, oShowArea);
		_this.drag(oActionBlock, oActionBar,minNumber, maxNumber, barLength, oShowArea);
	}
	
}

SlideBar.prototype = {
	//初始化(添加刻度线)
	addScale : function(slideBar, interval,min, total, barLength){
		// interval代表刻度之间间隔多少, total代表最大刻度
		// slideBar表示在哪个容器添加刻度

		var num = total / interval; //num为应该有多少个刻度
		var shu = null;
	},

	/* 在输入框输入值后自动滑动	*/
	autoSlide : function(actionBlock, actionBar,total, barLength, inputVal){
		//inputVal表示输入框中输入的值
		var _this = this;
		var target = (inputVal / total * barLength);
		_this.checkAndMove(actionBlock, actionBar, target);
	},

	/*	检查target(确认移动方向)并滑动	*/
	checkAndMove : function(actionBlock, actionBar, target){
		if(target > actionBar.offsetWidth){
			actionBarSpeed = 8;		//actionBar的移动度和方向
		}
		else if(target == actionBar.offsetWidth){
			return;
		}
		else if(target < actionBar.offsetWidth){
			actionBarSpeed = -8;
		}
		
		var timer = setInterval(function(){
			var actionBarPace = actionBar.offsetWidth + actionBarSpeed;

			if(Math.abs(actionBarPace - target) < 10){
				actionBarPace = target;
				clearInterval(timer);
			}
			actionBar.style.width = actionBarPace + 'px';
			actionBlock.style.left = actionBar.offsetWidth - (actionBlock.offsetWidth / 2) + 'px';
		},30);
	},

	/*	鼠标点击刻度滑动块自动滑动	*/
	clickSlide : function(actionBlock, actionBar, slideBar,min, total, barLength, showArea){
		var _this = this;
		slideBar.onclick = function(ev){
			var ev = ev || event;
			var target = ev.clientX - slideBar.offsetLeft;
			if(target <= 0){
				//表示鼠标已经超出那个范围
				target = min;
			}
			if(target > barLength){
				target = barLength;
			}
			_this.checkAndMove(actionBlock, actionBar, target);
			if(showArea){
				showArea.value = Math.ceil(target / barLength * total);	
			}
		}
	},

	/*	鼠标按着拖动滑动条	*/
	drag : function(actionBlock, actionBar,min, total, barLength, showArea){
		/*	参数分别是点击滑动的那个块,滑动的距离,滑动条的最大数值,显示数值的地方(输入框)	*/
		actionBlock.onmousedown = function(ev) {
			var ev = ev || event;
			var thisBlock = this;
			var disX = ev.clientX;
			var currentLeft = thisBlock.offsetLeft;

			document.onmousemove = function(ev) {
				var ev = ev || event;
				var left = ev.clientX - disX;

				if (currentLeft + left <= (barLength - thisBlock.offsetWidth / 2 ) && currentLeft + left >= 0 - thisBlock.offsetWidth / 2) {
					thisBlock.style.left = currentLeft + left + 'px';
					if(currentLeft + left + (actionBlock.offsetWidth / 2)) {
						actionBar.style.width = currentLeft + left + (actionBlock.offsetWidth / 2) + 'px';
					} else {
						actionBar.style.width = min + 'px';
					}
					if(showArea){
						showArea.value = Math.ceil(actionBar.offsetWidth / barLength * total);
					}
				}
				return false;
			}

			document.onmouseup = function() {
				document.onmousemove = document.onmouseup = null;
			}

			return false;
		}
	},

	getStyle : function(obj, attr){
		return obj.currentStyle?obj.currentStyle[attr]:getComputedStyle(obj)[attr];
	}
}