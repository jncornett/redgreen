$(function() {

  var Entry = Backbone.Model.extend({
    idAttribute: 'key',
    defaults: function() {
        return {
            'key': '',
            'ok': false,
            'data': null,
            'updated': null
        };
    },
    toObject: function() {
      var obj = this.toJSON();
      if (obj.updated) {
        obj.updated = moment(obj.updated).fromNow();
      }
      return obj;
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
      this.listenTo(this.model, 'remove', this.remove);
    },
    render: function() {
      this.$el.html(this.template(this.model.toObject()));
      this.$el.toggleClass("green", this.model.get('ok'));
      this.$el.addClass("redgreen-entry");
      return this;
    }
  });

  var AppView = Backbone.View.extend({
    el: $('#app'),
    initialize: function() {
      this.listenTo(Entries, 'add', this.addOne);
      this.listenTo(Entries, 'reset', this.addAll);
      this.listenTo(Entries, 'all', this.render);
      Entries.fetch({reset: true});
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
  setInterval(function() {
      Entries.fetch();
  }, 1000);
});
