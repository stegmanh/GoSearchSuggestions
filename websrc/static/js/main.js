"use strict";

var demo = new Vue({
    el: '#app',
    data: {
        term: null,
        results: null,
        selected: null
    },

    methods: {
    	autoComplete: function(e) {
    		switch (e.keyCode) {
    			case 38:
    				console.log(38)
    				return
    			case 40:
    				console.log(40)
    				return
    		}
			var self = this
			if (self.term.length < 1) {
				self.results = null
				return
			}
    		var xhr = new XMLHttpRequest()
	      xhr.open('GET', "/autocomplete?q=" + self.term)
	      xhr.onload = function () {
        		self.results = JSON.parse(xhr.responseText)
	      }
	      xhr.send()
    	}
    }
})

Vue.filter('take', function(value, limit) {
	return value.slice(0, limit)
})