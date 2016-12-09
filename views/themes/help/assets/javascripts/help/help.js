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
                .on(EVENT_CLICK, '.qor-help__lists a', this.loadDoc)
                .on(EVENT_KEYUP, '.qor-help__search', this.search);
        },

        unbind: function() {
            this.$element
                .off(EVENT_CLICK, '.qor-help__lists a', this.loadDoc)
                .off(EVENT_KEYUP, '.qor-help__search', this.search);
        },

        search: function(e) {
            if (e.keyCode == 13) {
                var $category = $('.qor-help__search-category'),
                    $input = $('.qor-help__search'),
                    $list = $('.qor-help__lists'),
                    $loading = $(QorHelpDocument.TEMPLATE_LOADING),
                    url = [
                        $(this).data().helpFilterUrl,
                        '?',
                        $category.prop('name'),
                        '=',
                        $category.find('option:selected').text(),
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
                        $list.html($(html).find('.qor-help__lists').html()).show();
                        $loading.remove();
                    },
                    error: function(xhr, textStatus, errorThrown) {
                        $list.show();
                        $loading.remove();
                        window.alert([textStatus, errorThrown].join(': '));
                    }
                });

            }
        },

        loadDoc: function(e) {
            var $element = $(e.target),
                $li = $element.closest('li'),
                $preview = $li.find('.qor-doc__preview'),
                $loading = $(QorHelpDocument.TEMPLATE_LOADING),
                url = $element.data().inlineUrl;

            if ($preview.size()) {
                $preview.toggle();
                return false;
            }

            $.ajax(url, {
                method: 'GET',
                dataType: 'html',
                processData: false,
                contentType: false,
                beforeSend: function() {
                    $li.append($loading);
                    window.componentHandler.upgradeElement($loading.children()[0]);
                },
                success: function(html) {
                    $($(html).find('.qor-form-container').html()).appendTo($li).addClass('qor-fieldset qor-doc__preview');
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