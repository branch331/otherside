<template>
    <b-navbar toggleable="lg" type="dark" variant="info">
      <b-navbar-brand href="#">Otherside</b-navbar-brand>

      <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>

      <b-collapse id="nav-collapse" is-nav>

        <!-- Right aligned nav items -->

        <b-navbar-nav class="ml-auto">
          <b-nav-item href="#">Link</b-nav-item>
          <b-nav-item href="#">Disabled</b-nav-item>

          <b-nav-item-dropdown right>
            <!-- Using 'button-content' slot -->
            <template v-slot:button-content>
              My Account
            </template>
            <b-dropdown-item href="#">Profile</b-dropdown-item>
            <b-dropdown-item href="#">Sign Out</b-dropdown-item>
          </b-nav-item-dropdown>
        </b-navbar-nav>
      </b-collapse>
    </b-navbar>
</template>

<script>
export default {
    name: 'AppHeader'
}
</script>
