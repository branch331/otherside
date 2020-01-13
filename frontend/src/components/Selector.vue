<template>
    <div class="selections">
        <h5>{{index}}. Pick some <strong><em>{{selectorName}}</em></strong>...</h5>

        <input type="text" 
            class="searchBar"
            v-model="searchTerm" 
            @focus="showSearchResults = true"
            @blur="searchBarRemovedFocus()">
           
        <div class="searchResults" v-if="showSearchResults">
            <div class = "searchResultItem" v-for="item in filteredItems" v-bind:key="item" v-on:mousedown="transferArrayValue(availableItems, selectedItems, item)">
                {{item}}
            </div>
        </div>

        <ul class="selectedItemList" v-if="selectedItems">
            <li class="selectedItem" v-for="item in selectedItems" v-bind:key="item" v-on:click="transferArrayValue(selectedItems, availableItems, item)">
                {{item}} <span class="selectedX">X</span>
            </li>
        </ul>
    </div>
</template>

<script>
export default {
    name: 'Selector',
    props: {
        selectorName: String,
        selectedItems: Array,
        availableItems: Array,
        maxItems: String,
        index: String
    },
    data () {
        return {
            searchTerm: null,
            showSearchResults: false
        }
    },
    computed: {
        filteredItems: {
            get: function() {
                var filteredArray = new Array();

                if (this.searchTerm == null) {
                    return this.availableItems;
                }

                this.availableItems.forEach(element => {
                    if (element.toLowerCase().includes(this.searchTerm.toLowerCase())) {
                        filteredArray.push(element);
                    }
                });

                return filteredArray;
            }
        }
    },
    methods: {
        transferArrayValue: function(sourceArray, targetArray, value) {
            var index = sourceArray.indexOf(value);
            if (index > -1) {
                if (targetArray == this.selectedItems && targetArray.length == this.maxItems) {
                    alert(`Please select less than ${this.maxItems} items.`);
                } else {
                    sourceArray.splice(index, 1);
                    targetArray.push(value)
                }
            }
        },
        searchBarRemovedFocus: function() {
            this.showSearchResults = false;
            this.searchTerm = null;
        }
    }
}
</script>