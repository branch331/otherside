import { shallowMount } from '@vue/test-utils'
import Selector from '../src/components/Selector.vue'
import Vue from 'vue'

describe('Selector', () => {
    let wrapper;

    beforeEach(() => {
        wrapper = shallowMount(Selector, {
            propsData: {
                selectorName: "testSelector",
                selectedItems: ['Sin City', 'Atlantis'],
                availableItems: ['New Delhi', 'Loverstown', 'Penicillin'],
                maxItems: "4"
            }
        })
    })

    it('populates selected items properly', () => {
        expect(wrapper.vm.selectedItems.length).toBe(2)

        expect(wrapper.text()).toContain('Sin City')
        expect(wrapper.text()).toContain('Atlantis')
        expect(wrapper.text()).not.toContain('New Delhi')
    })

    it('populates available items properly', () => {
        expect(wrapper.vm.availableItems.length).toBe(3)
    })

    it('toggles showSearchResults when search bar is focused or blurred', () => {
        const searchBar = wrapper.find('.searchBar')

        expect(wrapper.vm.showSearchResults).toBe(false)

        searchBar.trigger('focus')
        expect(wrapper.vm.showSearchResults).toBe(true)

        searchBar.trigger('blur')
        expect(wrapper.vm.showSearchResults).toBe(false)   
    })

    it('filters available items properly', () => {
        const searchBar = wrapper.find('.searchBar')
        searchBar.trigger('focus')

        searchBar.element.value = 'lover'
        searchBar.trigger('input')
        expect(wrapper.vm.filteredItems.length).toBe(1)

        searchBar.element.value = 'i'
        searchBar.trigger('input')
        expect(wrapper.vm.filteredItems.length).toBe(2)
    })

    it('transfers values from selected items to available items', () => {
        const selectedItem = wrapper.find('.selectedItem')
        selectedItem.trigger('click')
        expect(wrapper.vm.selectedItems.length).toBe(1)
        expect(wrapper.vm.availableItems.length).toBe(4)
    })

    it('transfers values from available items to selected items', async () => {
        const searchBar = wrapper.find('.searchBar')
        searchBar.trigger('focus')

        searchBar.element.value = 'lover'
        searchBar.trigger('input')

        expect(wrapper.vm.availableItems.length).toBe(3)

        await Vue.nextTick()
        const availableItem = wrapper.find('.searchResultItem')
        availableItem.trigger('mousedown')

        expect(wrapper.vm.availableItems.length).toBe(2)
    })

    it('prevents adding more selectedItems than the max', async () => {
        const searchBar = wrapper.find('.searchBar')
        searchBar.trigger('focus')
        searchBar.element.value = ''
        searchBar.trigger('input')
        
        expect(wrapper.vm.selectedItems.length).toBe(2)

        await Vue.nextTick()
        const availableItems = wrapper.findAll('.searchResultItem')

        availableItems.at(0).trigger('mousedown')
        await Vue.nextTick()
        expect(wrapper.vm.selectedItems.length).toBe(3)

        availableItems.at(1).trigger('mousedown')
        await Vue.nextTick()
        expect(wrapper.vm.selectedItems.length).toBe(4)

        availableItems.at(2).trigger('mousedown')
        await Vue.nextTick()
        expect(wrapper.vm.selectedItems.length).toBe(4)
    })
})