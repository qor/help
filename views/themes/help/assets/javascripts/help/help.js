(function(factory) {
    if (typeof define === 'function' && define.amd) {
        // AMD. Register as anonymous module.
        define(['jquery'], factory);
    } else if (typeof exports === 'object') {
        // Node / CommonJS
        factory(require('jquery'));
    } else {
        // Browser globals.
        factory(jQuery);
    }
})(function($) {

    'use strict';

    var NAMESPACE = 'qor.help';
    var EVENT_ENABLE = 'enable.' + NAMESPACE;
    var EVENT_DISABLE = 'disable.' + NAMESPACE;
    var EVENT_CLICK = 'click.' + NAMESPACE;
    var EVENT_KEYUP = 'keyup.' + NAMESPACE;
    var EVENT_CHANGE = 'change.' + NAMESPACE;


    function QorHelpDocument(element, options) {
        this.$element = $(element);
        this.options = $.extend({}, QorHelpDocument.DEFAULTS, $.isPlainObject(options) && options);
        this.init();
    }

    QorHelpDocument.prototype = {
        constructor: QorHelpDocument,

        init: function() {
            this.bind();
        },

        bind: function() {
            this.$element
                .on(EVENT_CLICK, '.qor-help__lists [data-inline-url]', this.loadDoc)
                .on(EVENT_KEYUP, '.qor-help__search', this.searchKeyup.bind(this))
                .on(EVENT_CLICK, '.qor-help__search-button', this.search.bind(this))
                .on(EVENT_CHANGE, '.qor-help__search-category', this.search.bind(this))
        },

        unbind: function() {
            this.$element
                .off(EVENT_CLICK, '.qor-help__lists [data-inline-url]', this.loadDoc)
                .off(EVENT_KEYUP, '.qor-help__search', this.searchKeyup.bind(this))
                .off(EVENT_CLICK, '.qor-help__search-button', this.search.bind(this))
                .off(EVENT_CHANGE, '.qor-help__search-category', this.search.bind(this))
        },

        searchKeyup: function(e) {
            if (e.keyCode == 13) {
                this.searchAction();
            }
        },

        search: function() {
            this.searchAction();
        },

        searchAction: function() {
            var $category = $('.qor-help__search-category'),
                $input = $('.qor-help__search'),
                $list = $('.qor-help__body'),
                $loading = $(QorHelpDocument.TEMPLATE_LOADING),
                url = [
                    $input.data().helpFilterUrl,
                    '?',
                    $category.prop('name'),
                    '=',
                    $category.val(),
                    '&',
                    $input.prop('name'),
                    '=',
                    $input.val()
                ].join('');

            $.ajax(url, {
                method: 'GET',
                dataType: 'html',
                processData: false,
                contentType: false,
                beforeSend: function() {
                    $list.hide().after($loading);
                    window.componentHandler.upgradeElement($loading.children()[0]);
                },

                success: function(html) {
                    $(".qor-slideout__title .qor-doc__close").remove();
                    $list.html($(html).find('.qor-help__body').html()).show();
                    $loading.remove();
                },
                error: function(xhr, textStatus, errorThrown) {
                    $list.show();
                    $loading.remove();
                    window.alert([textStatus, errorThrown].join(': '));
                }
            });
        },

        loadDoc: function(e) {
            var $this = this,
                $element = $(e.target),
                $index = $(".qor-help__index"),
                $loading = $(QorHelpDocument.TEMPLATE_LOADING),
                $help_body = $('.qor-help__body'),
                url = $element.data().inlineUrl;

            $index.hide();

            $.ajax(url, {
                method: 'GET',
                dataType: 'html',
                processData: false,
                contentType: false,
                beforeSend: function() {
                    $help_body.append($loading);
                    window.componentHandler.upgradeElement($loading.children()[0]);
                },
                success: function(html) {
                    $(html).find('.qor-page__show').appendTo($help_body).addClass('qor-doc__preview');
                    $(QorHelpDocument.TEMPLATE_PREVIEW_CLOSE).prependTo($(".qor-slideout__title"));
                    $(".qor-slideout__title .qor-doc__close").click(function() {
                        $index.show();
                        $('.qor-doc__preview').hide();
                        $(".qor-slideout__title .qor-doc__close").remove();
                    });
                    $loading.remove();
                },
                error: function(xhr, textStatus, errorThrown) {
                    $loading.remove();
                    window.alert([textStatus, errorThrown].join(': '));
                }
            });

            return false;
        },

        destroy: function() {
            this.unbind();
            this.$element.removeData(NAMESPACE);
        }
    };

    QorHelpDocument.TEMPLATE_LOADING = '<div style="text-align: center; margin-top: 30px;"><div class="mdl-spinner mdl-js-spinner is-active qor-layout__bottomsheet-spinner"></div></div>';
    QorHelpDocument.TEMPLATE_PREVIEW_CLOSE = '<a href="javascript://" class="qor-doc__close"><i class="material-icons">keyboard_backspace</i></a>';

    QorHelpDocument.plugin = function(options) {
        return this.each(function() {
            var $this = $(this);
            var data = $this.data(NAMESPACE);
            var fn;

            if (!data) {
                if (/destroy/.test(options)) {
                    return;
                }
                $this.data(NAMESPACE, (data = new QorHelpDocument(this, options)));
            }

            if (typeof options === 'string' && $.isFunction(fn = data[options])) {
                fn.apply(data);
            }
        });
    };


    $(function() {
        var selector = '[data-toggle="qor.help"]';

        $(document).
        on(EVENT_DISABLE, function(e) {
            QorHelpDocument.plugin.call($(selector, e.target), 'destroy');
        }).
        on(EVENT_ENABLE, function(e) {
            QorHelpDocument.plugin.call($(selector, e.target));
        }).
        triggerHandler(EVENT_ENABLE);
    });

    return QorHelpDocument;
});
