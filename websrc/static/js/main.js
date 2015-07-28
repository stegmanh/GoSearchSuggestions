"use strict";

var demo = new Vue({
    el: '#app',
    data: {
        term: null,
        results: null,
        selected: -1
    },

    methods: {
    	autoComplete: function(e) {
 			var self = this
    		switch (e.keyCode) {
    			case 38:
    				return
    			case 40:
			    	if (self.selected > -1) {
    					self.selected--
    					return
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