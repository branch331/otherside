import { shallowMount } from '@vue/test-utils'
import Selector from '../src/components/Selector.vue'

describe('Selector', () => {
    const wrapper = shallowMount(Selector, {
        propsData: {
            selectorName: "testSelector",
            selectedItems: ['Sin City', 'Atlantis'],
            availableItems: ['New Delhi', 'Loverstown', 'Penicillin'],
            maxItems: "4"
        }
    })

    it('populates selected items properly', () => {
        expect(wrapper.text()).toContain('Sin City')
        expect(wrapper.text()).toContain('Atlantis')
        expect(wrapper.text()).not.toContain('New Delhi')
    })

    it('populates available items properly', () => {

    })

    it('prevents adding more selectedItems than the max', () => {

    })

    it('transfers values from selected items to available items', () => {

    })

    it('transfers values from available items to selected items', () => {

    })

    it('toggles showSearchResults when search bar is focused or blurred', () => {
        const searchBar = wrapper.find('.searchBar')

        expect(wrapper.vm.showSearchResults).toBe(false)

        searchBar.trigger('focus')
        expect(wrapper.vm.showSearchResults).toBe(true)

        searchBar.trigger('blur')
        expect(wrapper.vm.showSearchResults).toBe(false)   
    })

    it('filters search results based on entered values', () => {

    })
})