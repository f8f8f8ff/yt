// https://stackoverflow.com/a/35029300
var lastChecked = null;
var checkboxes = document.querySelectorAll('input[type="checkbox"]');

for (var i=0;i<checkboxes.length;i++){
	checkboxes[i].setAttribute('data-index',i);
}

for (var i=0;i<checkboxes.length;i++){
	checkboxes[i].addEventListener("click",function(e){

		if(lastChecked && e.shiftKey) {
			var i = parseInt(lastChecked.getAttribute('data-index'));
			var j = parseInt(this.getAttribute('data-index'));
			var check_or_uncheck = this.checked;

			var low = i; var high=j;
			if (i>j){
				var low = j; var high=i; 
			}

			for(var c=0;c<checkboxes.length;c++){
				if (low <= c && c <=high){
					checkboxes[c].checked = check_or_uncheck;
				}   
			}
		} 
		lastChecked = this;
	});
}
