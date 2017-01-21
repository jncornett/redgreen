$(function() {

  var Entry = Backbone.Model.extend({
    idAttribute: 'key',
    defaults: function() {
        return {
            'key': '',
            'ok': false,
            'data': [],
            'updated': 0
        };
    }
  });

  var EntryList = Backbone.Collection.extend({
    url: '/api',
    model: Entry
  });

  var Entries = new EntryList;

  var EntryView = Backbone.View.extend({
    tagName: 'div',
    template: _.template($('#entry-template').html()),
    initialize: function() {
      this.listenTo(this.model, 'change', this.render);
      this.listenTo(this.model, 'destroy', this.remove);
    },
    render: function() {
      this.$el.html(this.template(this.model.toJSON()));
      this.$el.addClass("entry");
      this.$el.toggleClass("green", this.model.get("ok"));
      return this;
    }
  });

  var AppView = Backbone.View.extend({
    el: $('#app'),
    initialize: function() {
      this.listenTo(Entries, 'add', this.addOne);
      this.listenTo(Entries, 'reset', this.addAll);
      this.listenTo(Entries, 'all', this.render);
      Entries.fetch();
    },
    addOne: function(entry) {
      var view = new EntryView({model: entry});
      this.$el.append(view.render().el);
    },
    addAll: function() {
      Entries.each(this.addOne, this);
    }
  });

  var App = new AppView;
  // refresh the collection periodically
  setInterval(function() { Entries.fetch(); }, 1000);
});
