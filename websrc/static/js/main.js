"use strict";

var demo = new Vue({
    el: '#app',
    data: {
        term: null,
        results: null,
        selected: -1,
        focus: false
    },

    methods: {
    	autoComplete: function(e) {
 			var self = this
 			self.focus = true
 			if (e.keyCode == 38 || e.keyCode == 40) {
				switch (e.keyCode) {
	    			case 38:
	    				self.selected--
	    				if (self.selected < 0) {
	    					self.selected = 9
	    				}
				    	break;
	    			case 40:
	    				self.selected++
	    				if (self.selected > 9) {
	    					self.selected = 0
	    				}
	    				break;
		   			}
		   			self.term = self.$children[self.selected].$el.innerHTML
	    			self.$children.forEach(function(vEl) {
	    				vEl.$el.removeAttribute("class", "selected")
	    			})
	    			self.$children[self.selected].$el.setAttribute("class", "selected")
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

    	changeDisplays: function() {
    		var self = this
    		//self.focus = false
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