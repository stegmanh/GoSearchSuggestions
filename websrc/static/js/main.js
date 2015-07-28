"use strict";

var demo = new Vue({
    el: '#app',
    data: {
        term: null,
        results: null,
        selected: -1
        prev: null
    },

    methods: {
    	autoComplete: function(e) {
 			var self = this
 			if (e.keyCode == 38 || e.keyCode == 40) {
				switch (e.keyCode) {
	    			case 38:
				    	self.selected--
				    	break;
	    			case 40:
	    				self.selected++
	    				break;
		   			}
		   			if (self.selected < -1) {
		   				self.selected = 9
		   			}
		   			console.log(self.selected)
		   			self.term = self.$children[self.selected % 10].$el.innerHTML
	    			self.$children[self.selected % 10].$el.setAttribute("class", "selected")
	    			if (self.selected > 0) {
		    			self.$children[(self.selected - 1) % 10].$el.removeAttribute("class", "selected")
	    			}
				return
 			}
			if (self.term.length < 1) {
				self.results = null
				return
			}
    		var xhr = new XMLHttpRequest()
	      xhr.open('GET', "/autocomplete?q=" + self.term)
	      xhr.onload = function () {
        		self.results = JSON.parse(xhr.responseText)
	 			self.selected = -1
	      }
	      xhr.send()
    	},

    	changeInput: function(el) {
    		var self = this
    		self.term = el.$el.innerHTML
    		self.results = null
    	}
    }
})

Vue.filter('take', function(value, limit) {
	return value.slice(0, limit)
})