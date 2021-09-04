package datepicker

import (
	"encoding/base64"
	"io"
	"strings"

	"github.com/evocert/kwe/resources"
)

const datepickercss string = `.qs-datepicker-container{font-size:1rem;font-family:sans-serif;color:#000;position:absolute;width:15.625em;display:-webkit-box;display:-ms-flexbox;display:flex;-webkit-box-orient:vertical;-webkit-box-direction:normal;-ms-flex-direction:column;flex-direction:column;z-index:9001;-webkit-user-select:none;-moz-user-select:none;-ms-user-select:none;user-select:none;border:1px solid grey;border-radius:.263921875em;overflow:hidden;background:#fff;-webkit-box-shadow:0 1.25em 1.25em -.9375em rgba(0,0,0,.3);box-shadow:0 1.25em 1.25em -.9375em rgba(0,0,0,.3)}.qs-datepicker-container *{-webkit-box-sizing:border-box;box-sizing:border-box}.qs-centered{position:fixed;top:50%;left:50%;-webkit-transform:translate(-50%,-50%);-ms-transform:translate(-50%,-50%);transform:translate(-50%,-50%)}.qs-hidden{display:none}.qs-overlay{position:absolute;top:0;left:0;background:rgba(0,0,0,.75);color:#fff;width:100%;height:100%;padding:.5em;z-index:1;opacity:1;-webkit-transition:opacity .3s;transition:opacity .3s;display:-webkit-box;display:-ms-flexbox;display:flex;-webkit-box-orient:vertical;-webkit-box-direction:normal;-ms-flex-direction:column;flex-direction:column}.qs-overlay.qs-hidden{opacity:0;z-index:-1}.qs-overlay .qs-overlay-year{background:rgba(0,0,0,0);border:none;border-bottom:1px solid #fff;border-radius:0;color:#fff;font-size:.875em;padding:.25em 0;width:80%;text-align:center;margin:0 auto;display:block}.qs-overlay .qs-overlay-year::-webkit-inner-spin-button{-webkit-appearance:none}.qs-overlay .qs-close{padding:.5em;cursor:pointer;position:absolute;top:0;right:0}.qs-overlay .qs-submit{border:1px solid #fff;border-radius:.263921875em;padding:.5em;margin:0 auto auto;cursor:pointer;background:hsla(0,0%,50.2%,.4)}.qs-overlay .qs-submit.qs-disabled{color:grey;border-color:grey;cursor:not-allowed}.qs-overlay .qs-overlay-month-container{display:-webkit-box;display:-ms-flexbox;display:flex;-ms-flex-wrap:wrap;flex-wrap:wrap;-webkit-box-flex:1;-ms-flex-positive:1;flex-grow:1}.qs-overlay .qs-overlay-month{display:-webkit-box;display:-ms-flexbox;display:flex;-webkit-box-pack:center;-ms-flex-pack:center;justify-content:center;-webkit-box-align:center;-ms-flex-align:center;align-items:center;width:calc(100% / 3);cursor:pointer;opacity:.5;-webkit-transition:opacity .15s;transition:opacity .15s}.qs-overlay .qs-overlay-month.active,.qs-overlay .qs-overlay-month:hover{opacity:1}.qs-controls{width:100%;display:-webkit-box;display:-ms-flexbox;display:flex;-webkit-box-pack:justify;-ms-flex-pack:justify;justify-content:space-between;-webkit-box-align:center;-ms-flex-align:center;align-items:center;-webkit-box-flex:1;-ms-flex-positive:1;flex-grow:1;-ms-flex-negative:0;flex-shrink:0;background:#d3d3d3;-webkit-filter:blur(0);filter:blur(0);-webkit-transition:-webkit-filter .3s;transition:-webkit-filter .3s;transition:filter .3s;transition:filter .3s,-webkit-filter .3s}.qs-controls.qs-blur{-webkit-filter:blur(5px);filter:blur(5px)}.qs-arrow{height:1.5625em;width:1.5625em;position:relative;cursor:pointer;border-radius:.263921875em;-webkit-transition:background .15s;transition:background .15s}.qs-arrow:hover{background:rgba(0,0,0,.1)}.qs-arrow:hover.qs-left:after{border-right-color:#000}.qs-arrow:hover.qs-right:after{border-left-color:#000}.qs-arrow:after{content:"";border:.390625em solid rgba(0,0,0,0);position:absolute;top:50%;-webkit-transition:border .2s;transition:border .2s}.qs-arrow.qs-left:after{border-right-color:grey;right:50%;-webkit-transform:translate(25%,-50%);-ms-transform:translate(25%,-50%);transform:translate(25%,-50%)}.qs-arrow.qs-right:after{border-left-color:grey;left:50%;-webkit-transform:translate(-25%,-50%);-ms-transform:translate(-25%,-50%);transform:translate(-25%,-50%)}.qs-month-year{font-weight:700;-webkit-transition:border .2s;transition:border .2s;border-bottom:1px solid rgba(0,0,0,0);cursor:pointer}.qs-month-year:hover{border-bottom:1px solid grey}.qs-month-year:active:focus,.qs-month-year:focus{outline:none}.qs-month{padding-right:.5ex}.qs-year{padding-left:.5ex}.qs-squares{display:-webkit-box;display:-ms-flexbox;display:flex;-ms-flex-wrap:wrap;flex-wrap:wrap;padding:.3125em;-webkit-filter:blur(0);filter:blur(0);-webkit-transition:-webkit-filter .3s;transition:-webkit-filter .3s;transition:filter .3s;transition:filter .3s,-webkit-filter .3s}.qs-squares.qs-blur{-webkit-filter:blur(5px);filter:blur(5px)}.qs-square{width:calc(100% / 7);height:1.5625em;display:-webkit-box;display:-ms-flexbox;display:flex;-webkit-box-align:center;-ms-flex-align:center;align-items:center;-webkit-box-pack:center;-ms-flex-pack:center;justify-content:center;cursor:pointer;-webkit-transition:background .1s;transition:background .1s;border-radius:.263921875em}.qs-square:not(.qs-empty):not(.qs-disabled):not(.qs-day):not(.qs-active):hover{background:orange}.qs-current{font-weight:700;text-decoration:underline}.qs-active,.qs-range-end,.qs-range-start{background:#add8e6}.qs-range-start:not(.qs-range-6){border-top-right-radius:0;border-bottom-right-radius:0}.qs-range-middle{background:#d4ebf2}.qs-range-middle:not(.qs-range-0):not(.qs-range-6){border-radius:0}.qs-range-middle.qs-range-0{border-top-right-radius:0;border-bottom-right-radius:0}.qs-range-end:not(.qs-range-0),.qs-range-middle.qs-range-6{border-top-left-radius:0;border-bottom-left-radius:0}.qs-disabled,.qs-outside-current-month{opacity:.2}.qs-disabled{cursor:not-allowed}.qs-day,.qs-empty{cursor:default}.qs-day{font-weight:700;color:grey}.qs-event{position:relative}.qs-event:after{content:"";position:absolute;width:.46875em;height:.46875em;border-radius:50%;background:#07f;bottom:0;right:0}`

const datepickerjs string = `IWZ1bmN0aW9uKGUsdCl7Im9iamVjdCI9PXR5cGVvZiBleHBvcnRzJiYib2JqZWN0Ij09dHlwZW9mIG1vZHVsZT9tb2R1bGUuZXhwb3J0cz10KCk6ImZ1bmN0aW9uIj09dHlwZW9mIGRlZmluZSYmZGVmaW5lLmFtZD9kZWZpbmUoW10sdCk6Im9iamVjdCI9PXR5cGVvZiBleHBvcnRzP2V4cG9ydHMuZGF0ZXBpY2tlcj10KCk6ZS5kYXRlcGlja2VyPXQoKX0od2luZG93LChmdW5jdGlvbigpe3JldHVybiBmdW5jdGlvbihlKXt2YXIgdD17fTtmdW5jdGlvbiBuKGEpe2lmKHRbYV0pcmV0dXJuIHRbYV0uZXhwb3J0czt2YXIgcj10W2FdPXtpOmEsbDohMSxleHBvcnRzOnt9fTtyZXR1cm4gZVthXS5jYWxsKHIuZXhwb3J0cyxyLHIuZXhwb3J0cyxuKSxyLmw9ITAsci5leHBvcnRzfXJldHVybiBuLm09ZSxuLmM9dCxuLmQ9ZnVuY3Rpb24oZSx0LGEpe24ubyhlLHQpfHxPYmplY3QuZGVmaW5lUHJvcGVydHkoZSx0LHtlbnVtZXJhYmxlOiEwLGdldDphfSl9LG4ucj1mdW5jdGlvbihlKXsidW5kZWZpbmVkIiE9dHlwZW9mIFN5bWJvbCYmU3ltYm9sLnRvU3RyaW5nVGFnJiZPYmplY3QuZGVmaW5lUHJvcGVydHkoZSxTeW1ib2wudG9TdHJpbmdUYWcse3ZhbHVlOiJNb2R1bGUifSksT2JqZWN0LmRlZmluZVByb3BlcnR5KGUsIl9fZXNNb2R1bGUiLHt2YWx1ZTohMH0pfSxuLnQ9ZnVuY3Rpb24oZSx0KXtpZigxJnQmJihlPW4oZSkpLDgmdClyZXR1cm4gZTtpZig0JnQmJiJvYmplY3QiPT10eXBlb2YgZSYmZSYmZS5fX2VzTW9kdWxlKXJldHVybiBlO3ZhciBhPU9iamVjdC5jcmVhdGUobnVsbCk7aWYobi5yKGEpLE9iamVjdC5kZWZpbmVQcm9wZXJ0eShhLCJkZWZhdWx0Iix7ZW51bWVyYWJsZTohMCx2YWx1ZTplfSksMiZ0JiYic3RyaW5nIiE9dHlwZW9mIGUpZm9yKHZhciByIGluIGUpbi5kKGEscixmdW5jdGlvbih0KXtyZXR1cm4gZVt0XX0uYmluZChudWxsLHIpKTtyZXR1cm4gYX0sbi5uPWZ1bmN0aW9uKGUpe3ZhciB0PWUmJmUuX19lc01vZHVsZT9mdW5jdGlvbigpe3JldHVybiBlLmRlZmF1bHR9OmZ1bmN0aW9uKCl7cmV0dXJuIGV9O3JldHVybiBuLmQodCwiYSIsdCksdH0sbi5vPWZ1bmN0aW9uKGUsdCl7cmV0dXJuIE9iamVjdC5wcm90b3R5cGUuaGFzT3duUHJvcGVydHkuY2FsbChlLHQpfSxuLnA9IiIsbihuLnM9MCl9KFtmdW5jdGlvbihlLHQsbil7InVzZSBzdHJpY3QiO24ucih0KTt2YXIgYT1bXSxyPVsiU3VuIiwiTW9uIiwiVHVlIiwiV2VkIiwiVGh1IiwiRnJpIiwiU2F0Il0saT1bIkphbnVhcnkiLCJGZWJydWFyeSIsIk1hcmNoIiwiQXByaWwiLCJNYXkiLCJKdW5lIiwiSnVseSIsIkF1Z3VzdCIsIlNlcHRlbWJlciIsIk9jdG9iZXIiLCJOb3ZlbWJlciIsIkRlY2VtYmVyIl0sbz17dDoidG9wIixyOiJyaWdodCIsYjoiYm90dG9tIixsOiJsZWZ0IixjOiJjZW50ZXJlZCJ9O2Z1bmN0aW9uIHMoKXt9dmFyIGw9WyJjbGljayIsImZvY3VzaW4iLCJrZXlkb3duIiwiaW5wdXQiXTtmdW5jdGlvbiBkKGUpe2wuZm9yRWFjaCgoZnVuY3Rpb24odCl7ZS5hZGRFdmVudExpc3RlbmVyKHQsZT09PWRvY3VtZW50P0w6WSl9KSl9ZnVuY3Rpb24gYyhlKXtyZXR1cm4gQXJyYXkuaXNBcnJheShlKT9lLm1hcChjKToiW29iamVjdCBPYmplY3RdIj09PXgoZSk/T2JqZWN0LmtleXMoZSkucmVkdWNlKChmdW5jdGlvbih0LG4pe3JldHVybiB0W25dPWMoZVtuXSksdH0pLHt9KTplfWZ1bmN0aW9uIHUoZSx0KXt2YXIgbj1lLmNhbGVuZGFyLnF1ZXJ5U2VsZWN0b3IoIi5xcy1vdmVybGF5IiksYT1uJiYhbi5jbGFzc0xpc3QuY29udGFpbnMoInFzLWhpZGRlbiIpO3Q9dHx8bmV3IERhdGUoZS5jdXJyZW50WWVhcixlLmN1cnJlbnRNb250aCksZS5jYWxlbmRhci5pbm5lckhUTUw9W2godCxlLGEpLGYodCxlLGEpLHYoZSxhKV0uam9pbigiIiksYSYmd2luZG93LnJlcXVlc3RBbmltYXRpb25GcmFtZSgoZnVuY3Rpb24oKXtNKCEwLGUpfSkpfWZ1bmN0aW9uIGgoZSx0LG4pe3JldHVyblsnPGRpdiBjbGFzcz0icXMtY29udHJvbHMnKyhuPyIgcXMtYmx1ciI6IiIpKyciPicsJzxkaXYgY2xhc3M9InFzLWFycm93IHFzLWxlZnQiPjwvZGl2PicsJzxkaXYgY2xhc3M9InFzLW1vbnRoLXllYXIiPicsJzxzcGFuIGNsYXNzPSJxcy1tb250aCI+Jyt0Lm1vbnRoc1tlLmdldE1vbnRoKCldKyI8L3NwYW4+IiwnPHNwYW4gY2xhc3M9InFzLXllYXIiPicrZS5nZXRGdWxsWWVhcigpKyI8L3NwYW4+IiwiPC9kaXY+IiwnPGRpdiBjbGFzcz0icXMtYXJyb3cgcXMtcmlnaHQiPjwvZGl2PicsIjwvZGl2PiJdLmpvaW4oIiIpfWZ1bmN0aW9uIGYoZSx0LG4pe3ZhciBhPXQuY3VycmVudE1vbnRoLHI9dC5jdXJyZW50WWVhcixpPXQuZGF0ZVNlbGVjdGVkLG89dC5tYXhEYXRlLHM9dC5taW5EYXRlLGw9dC5zaG93QWxsRGF0ZXMsZD10LmRheXMsYz10LmRpc2FibGVkRGF0ZXMsdT10LnN0YXJ0RGF5LGg9dC53ZWVrZW5kSW5kaWNlcyxmPXQuZXZlbnRzLHY9dC5nZXRSYW5nZT90LmdldFJhbmdlKCk6e30sbT0rdi5zdGFydCx5PSt2LmVuZCxwPWcobmV3IERhdGUoZSkuc2V0RGF0ZSgxKSksdz1wLmdldERheSgpLXUsRD13PDA/NzowO3Auc2V0TW9udGgocC5nZXRNb250aCgpKzEpLHAuc2V0RGF0ZSgwKTt2YXIgYj1wLmdldERhdGUoKSxxPVtdLFM9RCs3KigodytiKS83fDApO1MrPSh3K2IpJTc/NzowO2Zvcih2YXIgTT0xO008PVM7TSsrKXt2YXIgRT0oTS0xKSU3LHg9ZFtFXSxDPU0tKHc+PTA/dzo3K3cpLEw9bmV3IERhdGUocixhLEMpLFk9ZlsrTF0saj1DPDF8fEM+YixQPWo/QzwxPy0xOjE6MCxrPWomJiFsLE89az8iIjpMLmdldERhdGUoKSxOPStMPT0raSxfPUU9PT1oWzBdfHxFPT09aFsxXSxJPW0hPT15LEE9InFzLXNxdWFyZSAiK3g7WSYmIWsmJihBKz0iIHFzLWV2ZW50IiksaiYmKEErPSIgcXMtb3V0c2lkZS1jdXJyZW50LW1vbnRoIiksIWwmJmp8fChBKz0iIHFzLW51bSIpLE4mJihBKz0iIHFzLWFjdGl2ZSIpLChjWytMXXx8dC5kaXNhYmxlcihMKXx8XyYmdC5ub1dlZWtlbmRzfHxzJiYrTDwrc3x8byYmK0w+K28pJiYhayYmKEErPSIgcXMtZGlzYWJsZWQiKSwrZyhuZXcgRGF0ZSk9PStMJiYoQSs9IiBxcy1jdXJyZW50IiksK0w9PT1tJiZ5JiZJJiYoQSs9IiBxcy1yYW5nZS1zdGFydCIpLCtMPm0mJitMPHkmJihBKz0iIHFzLXJhbmdlLW1pZGRsZSIpLCtMPT09eSYmbSYmSSYmKEErPSIgcXMtcmFuZ2UtZW5kIiksayYmKEErPSIgcXMtZW1wdHkiLE89IiIpLHEucHVzaCgnPGRpdiBjbGFzcz0iJytBKyciIGRhdGEtZGlyZWN0aW9uPSInK1ArJyI+JytPKyI8L2Rpdj4iKX12YXIgUj1kLm1hcCgoZnVuY3Rpb24oZSl7cmV0dXJuJzxkaXYgY2xhc3M9InFzLXNxdWFyZSBxcy1kYXkiPicrZSsiPC9kaXY+In0pKS5jb25jYXQocSk7cmV0dXJuIFIudW5zaGlmdCgnPGRpdiBjbGFzcz0icXMtc3F1YXJlcycrKG4/IiBxcy1ibHVyIjoiIikrJyI+JyksUi5wdXNoKCI8L2Rpdj4iKSxSLmpvaW4oIiIpfWZ1bmN0aW9uIHYoZSx0KXt2YXIgbj1lLm92ZXJsYXlQbGFjZWhvbGRlcixhPWUub3ZlcmxheUJ1dHRvbjtyZXR1cm5bJzxkaXYgY2xhc3M9InFzLW92ZXJsYXknKyh0PyIiOiIgcXMtaGlkZGVuIikrJyI+JywiPGRpdj4iLCc8aW5wdXQgY2xhc3M9InFzLW92ZXJsYXkteWVhciIgcGxhY2Vob2xkZXI9IicrbisnIiBpbnB1dG1vZGU9Im51bWVyaWMiIC8+JywnPGRpdiBjbGFzcz0icXMtY2xvc2UiPiYjMTAwMDU7PC9kaXY+JywiPC9kaXY+IiwnPGRpdiBjbGFzcz0icXMtb3ZlcmxheS1tb250aC1jb250YWluZXIiPicrZS5vdmVybGF5TW9udGhzLm1hcCgoZnVuY3Rpb24oZSx0KXtyZXR1cm4nPGRpdiBjbGFzcz0icXMtb3ZlcmxheS1tb250aCIgZGF0YS1tb250aC1udW09IicrdCsnIj4nK2UrIjwvZGl2PiJ9KSkuam9pbigiIikrIjwvZGl2PiIsJzxkaXYgY2xhc3M9InFzLXN1Ym1pdCBxcy1kaXNhYmxlZCI+JythKyI8L2Rpdj4iLCI8L2Rpdj4iXS5qb2luKCIiKX1mdW5jdGlvbiBtKGUsdCxuKXt2YXIgYT10LmVsLHI9dC5jYWxlbmRhci5xdWVyeVNlbGVjdG9yKCIucXMtYWN0aXZlIiksaT1lLnRleHRDb250ZW50LG89dC5zaWJsaW5nOyhhLmRpc2FibGVkfHxhLnJlYWRPbmx5KSYmdC5yZXNwZWN0RGlzYWJsZWRSZWFkT25seXx8KHQuZGF0ZVNlbGVjdGVkPW4/dm9pZCAwOm5ldyBEYXRlKHQuY3VycmVudFllYXIsdC5jdXJyZW50TW9udGgsaSksciYmci5jbGFzc0xpc3QucmVtb3ZlKCJxcy1hY3RpdmUiKSxufHxlLmNsYXNzTGlzdC5hZGQoInFzLWFjdGl2ZSIpLHAoYSx0LG4pLG58fHEodCksbyYmKHkoe2luc3RhbmNlOnQsZGVzZWxlY3Q6bn0pLHQuZmlyc3QmJiFvLmRhdGVTZWxlY3RlZCYmKG8uY3VycmVudFllYXI9dC5jdXJyZW50WWVhcixvLmN1cnJlbnRNb250aD10LmN1cnJlbnRNb250aCxvLmN1cnJlbnRNb250aE5hbWU9dC5jdXJyZW50TW9udGhOYW1lKSx1KHQpLHUobykpLHQub25TZWxlY3QodCxuP3ZvaWQgMDpuZXcgRGF0ZSh0LmRhdGVTZWxlY3RlZCkpKX1mdW5jdGlvbiB5KGUpe3ZhciB0PWUuaW5zdGFuY2UuZmlyc3Q/ZS5pbnN0YW5jZTplLmluc3RhbmNlLnNpYmxpbmcsbj10LnNpYmxpbmc7dD09PWUuaW5zdGFuY2U/ZS5kZXNlbGVjdD8odC5taW5EYXRlPXQub3JpZ2luYWxNaW5EYXRlLG4ubWluRGF0ZT1uLm9yaWdpbmFsTWluRGF0ZSk6bi5taW5EYXRlPXQuZGF0ZVNlbGVjdGVkOmUuZGVzZWxlY3Q/KG4ubWF4RGF0ZT1uLm9yaWdpbmFsTWF4RGF0ZSx0Lm1heERhdGU9dC5vcmlnaW5hbE1heERhdGUpOnQubWF4RGF0ZT1uLmRhdGVTZWxlY3RlZH1mdW5jdGlvbiBwKGUsdCxuKXtpZighdC5ub25JbnB1dClyZXR1cm4gbj9lLnZhbHVlPSIiOnQuZm9ybWF0dGVyIT09cz90LmZvcm1hdHRlcihlLHQuZGF0ZVNlbGVjdGVkLHQpOnZvaWQoZS52YWx1ZT10LmRhdGVTZWxlY3RlZC50b0RhdGVTdHJpbmcoKSl9ZnVuY3Rpb24gdyhlLHQsbixhKXtufHxhPyhuJiYodC5jdXJyZW50WWVhcj0rbiksYSYmKHQuY3VycmVudE1vbnRoPSthKSk6KHQuY3VycmVudE1vbnRoKz1lLmNvbnRhaW5zKCJxcy1yaWdodCIpPzE6LTEsMTI9PT10LmN1cnJlbnRNb250aD8odC5jdXJyZW50TW9udGg9MCx0LmN1cnJlbnRZZWFyKyspOi0xPT09dC5jdXJyZW50TW9udGgmJih0LmN1cnJlbnRNb250aD0xMSx0LmN1cnJlbnRZZWFyLS0pKSx0LmN1cnJlbnRNb250aE5hbWU9dC5tb250aHNbdC5jdXJyZW50TW9udGhdLHUodCksdC5vbk1vbnRoQ2hhbmdlKHQpfWZ1bmN0aW9uIEQoZSl7aWYoIWUubm9Qb3NpdGlvbil7dmFyIHQ9ZS5wb3NpdGlvbi50b3Asbj1lLnBvc2l0aW9uLnJpZ2h0O2lmKGUucG9zaXRpb24uY2VudGVyZWQpcmV0dXJuIGUuY2FsZW5kYXJDb250YWluZXIuY2xhc3NMaXN0LmFkZCgicXMtY2VudGVyZWQiKTt2YXIgYT1lLnBvc2l0aW9uZWRFbC5nZXRCb3VuZGluZ0NsaWVudFJlY3QoKSxyPWUuZWwuZ2V0Qm91bmRpbmdDbGllbnRSZWN0KCksaT1lLmNhbGVuZGFyQ29udGFpbmVyLmdldEJvdW5kaW5nQ2xpZW50UmVjdCgpLG89ci50b3AtYS50b3ArKHQ/LTEqaS5oZWlnaHQ6ci5oZWlnaHQpKyJweCIscz1yLmxlZnQtYS5sZWZ0KyhuP3Iud2lkdGgtaS53aWR0aDowKSsicHgiO2UuY2FsZW5kYXJDb250YWluZXIuc3R5bGUuc2V0UHJvcGVydHkoInRvcCIsbyksZS5jYWxlbmRhckNvbnRhaW5lci5zdHlsZS5zZXRQcm9wZXJ0eSgibGVmdCIscyl9fWZ1bmN0aW9uIGIoZSl7cmV0dXJuIltvYmplY3QgRGF0ZV0iPT09eChlKSYmIkludmFsaWQgRGF0ZSIhPT1lLnRvU3RyaW5nKCl9ZnVuY3Rpb24gZyhlKXtpZihiKGUpfHwibnVtYmVyIj09dHlwZW9mIGUmJiFpc05hTihlKSl7dmFyIHQ9bmV3IERhdGUoK2UpO3JldHVybiBuZXcgRGF0ZSh0LmdldEZ1bGxZZWFyKCksdC5nZXRNb250aCgpLHQuZ2V0RGF0ZSgpKX19ZnVuY3Rpb24gcShlKXtlLmRpc2FibGVkfHwhZS5jYWxlbmRhckNvbnRhaW5lci5jbGFzc0xpc3QuY29udGFpbnMoInFzLWhpZGRlbiIpJiYhZS5hbHdheXNTaG93JiYoIm92ZXJsYXkiIT09ZS5kZWZhdWx0VmlldyYmTSghMCxlKSxlLmNhbGVuZGFyQ29udGFpbmVyLmNsYXNzTGlzdC5hZGQoInFzLWhpZGRlbiIpLGUub25IaWRlKGUpKX1mdW5jdGlvbiBTKGUpe2UuZGlzYWJsZWR8fChlLmNhbGVuZGFyQ29udGFpbmVyLmNsYXNzTGlzdC5yZW1vdmUoInFzLWhpZGRlbiIpLCJvdmVybGF5Ij09PWUuZGVmYXVsdFZpZXcmJk0oITEsZSksRChlKSxlLm9uU2hvdyhlKSl9ZnVuY3Rpb24gTShlLHQpe3ZhciBuPXQuY2FsZW5kYXIsYT1uLnF1ZXJ5U2VsZWN0b3IoIi5xcy1vdmVybGF5Iikscj1hLnF1ZXJ5U2VsZWN0b3IoIi5xcy1vdmVybGF5LXllYXIiKSxpPW4ucXVlcnlTZWxlY3RvcigiLnFzLWNvbnRyb2xzIiksbz1uLnF1ZXJ5U2VsZWN0b3IoIi5xcy1zcXVhcmVzIik7ZT8oYS5jbGFzc0xpc3QuYWRkKCJxcy1oaWRkZW4iKSxpLmNsYXNzTGlzdC5yZW1vdmUoInFzLWJsdXIiKSxvLmNsYXNzTGlzdC5yZW1vdmUoInFzLWJsdXIiKSxyLnZhbHVlPSIiKTooYS5jbGFzc0xpc3QucmVtb3ZlKCJxcy1oaWRkZW4iKSxpLmNsYXNzTGlzdC5hZGQoInFzLWJsdXIiKSxvLmNsYXNzTGlzdC5hZGQoInFzLWJsdXIiKSxyLmZvY3VzKCkpfWZ1bmN0aW9uIEUoZSx0LG4sYSl7dmFyIHI9aXNOYU4oKyhuZXcgRGF0ZSkuc2V0RnVsbFllYXIodC52YWx1ZXx8dm9pZCAwKSksaT1yP251bGw6dC52YWx1ZTtpZigxMz09PWUud2hpY2h8fDEzPT09ZS5rZXlDb2RlfHwiY2xpY2siPT09ZS50eXBlKWE/dyhudWxsLG4saSxhKTpyfHx0LmNsYXNzTGlzdC5jb250YWlucygicXMtZGlzYWJsZWQiKXx8dyhudWxsLG4saSk7ZWxzZSBpZihuLmNhbGVuZGFyLmNvbnRhaW5zKHQpKXtuLmNhbGVuZGFyLnF1ZXJ5U2VsZWN0b3IoIi5xcy1zdWJtaXQiKS5jbGFzc0xpc3Rbcj8iYWRkIjoicmVtb3ZlIl0oInFzLWRpc2FibGVkIil9fWZ1bmN0aW9uIHgoZSl7cmV0dXJue30udG9TdHJpbmcuY2FsbChlKX1mdW5jdGlvbiBDKGUpe2EuZm9yRWFjaCgoZnVuY3Rpb24odCl7dCE9PWUmJnEodCl9KSl9ZnVuY3Rpb24gTChlKXtpZighZS5fX3FzX3NoYWRvd19kb20pe3ZhciB0PWUud2hpY2h8fGUua2V5Q29kZSxuPWUudHlwZSxyPWUudGFyZ2V0LG89ci5jbGFzc0xpc3Qscz1hLmZpbHRlcigoZnVuY3Rpb24oZSl7cmV0dXJuIGUuY2FsZW5kYXIuY29udGFpbnMocil8fGUuZWw9PT1yfSkpWzBdLGw9cyYmcy5jYWxlbmRhci5jb250YWlucyhyKTtpZighKHMmJnMuaXNNb2JpbGUmJnMuZGlzYWJsZU1vYmlsZSkpaWYoImNsaWNrIj09PW4pe2lmKCFzKXJldHVybiBhLmZvckVhY2gocSk7aWYocy5kaXNhYmxlZClyZXR1cm47dmFyIGQ9cy5jYWxlbmRhcixjPXMuY2FsZW5kYXJDb250YWluZXIsaD1zLmRpc2FibGVZZWFyT3ZlcmxheSxmPXMubm9uSW5wdXQsdj1kLnF1ZXJ5U2VsZWN0b3IoIi5xcy1vdmVybGF5LXllYXIiKSx5PSEhZC5xdWVyeVNlbGVjdG9yKCIucXMtaGlkZGVuIikscD1kLnF1ZXJ5U2VsZWN0b3IoIi5xcy1tb250aC15ZWFyIikuY29udGFpbnMociksRD1yLmRhdGFzZXQubW9udGhOdW07aWYocy5ub1Bvc2l0aW9uJiYhbCkoYy5jbGFzc0xpc3QuY29udGFpbnMoInFzLWhpZGRlbiIpP1M6cSkocyk7ZWxzZSBpZihvLmNvbnRhaW5zKCJxcy1hcnJvdyIpKXcobyxzKTtlbHNlIGlmKHB8fG8uY29udGFpbnMoInFzLWNsb3NlIikpaHx8TSgheSxzKTtlbHNlIGlmKEQpRShlLHYscyxEKTtlbHNle2lmKG8uY29udGFpbnMoInFzLWRpc2FibGVkIikpcmV0dXJuO2lmKG8uY29udGFpbnMoInFzLW51bSIpKXt2YXIgYj1yLnRleHRDb250ZW50LGc9K3IuZGF0YXNldC5kaXJlY3Rpb24seD1uZXcgRGF0ZShzLmN1cnJlbnRZZWFyLHMuY3VycmVudE1vbnRoK2csYik7aWYoZyl7cy5jdXJyZW50WWVhcj14LmdldEZ1bGxZZWFyKCkscy5jdXJyZW50TW9udGg9eC5nZXRNb250aCgpLHMuY3VycmVudE1vbnRoTmFtZT1pW3MuY3VycmVudE1vbnRoXSx1KHMpO2Zvcih2YXIgTCxZPXMuY2FsZW5kYXIucXVlcnlTZWxlY3RvckFsbCgnW2RhdGEtZGlyZWN0aW9uPSIwIl0nKSxqPTA7IUw7KXt2YXIgUD1ZW2pdO1AudGV4dENvbnRlbnQ9PT1iJiYoTD1QKSxqKyt9cj1MfXJldHVybiB2b2lkKCt4PT0rcy5kYXRlU2VsZWN0ZWQ/bShyLHMsITApOnIuY2xhc3NMaXN0LmNvbnRhaW5zKCJxcy1kaXNhYmxlZCIpfHxtKHIscykpfW8uY29udGFpbnMoInFzLXN1Ym1pdCIpP0UoZSx2LHMpOmYmJnI9PT1zLmVsJiYoUyhzKSxDKHMpKX19ZWxzZSBpZigiZm9jdXNpbiI9PT1uJiZzKVMocyksQyhzKTtlbHNlIGlmKCJrZXlkb3duIj09PW4mJjk9PT10JiZzKXEocyk7ZWxzZSBpZigia2V5ZG93biI9PT1uJiZzJiYhcy5kaXNhYmxlZCl7dmFyIGs9IXMuY2FsZW5kYXIucXVlcnlTZWxlY3RvcigiLnFzLW92ZXJsYXkiKS5jbGFzc0xpc3QuY29udGFpbnMoInFzLWhpZGRlbiIpOzEzPT09dCYmayYmbD9FKGUscixzKToyNz09PXQmJmsmJmwmJk0oITAscyl9ZWxzZSBpZigiaW5wdXQiPT09bil7aWYoIXN8fCFzLmNhbGVuZGFyLmNvbnRhaW5zKHIpKXJldHVybjt2YXIgTz1zLmNhbGVuZGFyLnF1ZXJ5U2VsZWN0b3IoIi5xcy1zdWJtaXQiKSxOPXIudmFsdWUuc3BsaXQoIiIpLnJlZHVjZSgoZnVuY3Rpb24oZSx0KXtyZXR1cm4gZXx8IjAiIT09dD9lKyh0Lm1hdGNoKC9bMC05XS8pP3Q6IiIpOiIifSksIiIpLnNsaWNlKDAsNCk7ci52YWx1ZT1OLE8uY2xhc3NMaXN0WzQ9PT1OLmxlbmd0aD8icmVtb3ZlIjoiYWRkIl0oInFzLWRpc2FibGVkIil9fX1mdW5jdGlvbiBZKGUpe0woZSksZS5fX3FzX3NoYWRvd19kb209ITB9ZnVuY3Rpb24gaihlLHQpe2wuZm9yRWFjaCgoZnVuY3Rpb24obil7ZS5yZW1vdmVFdmVudExpc3RlbmVyKG4sdCl9KSl9ZnVuY3Rpb24gUCgpe1ModGhpcyl9ZnVuY3Rpb24gaygpe3EodGhpcyl9ZnVuY3Rpb24gTyhlLHQpe3ZhciBuPWcoZSksYT10aGlzLmN1cnJlbnRZZWFyLHI9dGhpcy5jdXJyZW50TW9udGgsaT10aGlzLnNpYmxpbmc7aWYobnVsbD09ZSlyZXR1cm4gdGhpcy5kYXRlU2VsZWN0ZWQ9dm9pZCAwLHAodGhpcy5lbCx0aGlzLCEwKSxpJiYoeSh7aW5zdGFuY2U6dGhpcyxkZXNlbGVjdDohMH0pLHUoaSkpLHUodGhpcyksdGhpcztpZighYihlKSl0aHJvdyBuZXcgRXJyb3IoImBzZXREYXRlYCBuZWVkcyBhIEphdmFTY3JpcHQgRGF0ZSBvYmplY3QuIik7aWYodGhpcy5kaXNhYmxlZERhdGVzWytuXXx8bjx0aGlzLm1pbkRhdGV8fG4+dGhpcy5tYXhEYXRlKXRocm93IG5ldyBFcnJvcigiWW91IGNhbid0IG1hbnVhbGx5IHNldCBhIGRhdGUgdGhhdCdzIGRpc2FibGVkLiIpO3RoaXMuZGF0ZVNlbGVjdGVkPW4sdCYmKHRoaXMuY3VycmVudFllYXI9bi5nZXRGdWxsWWVhcigpLHRoaXMuY3VycmVudE1vbnRoPW4uZ2V0TW9udGgoKSx0aGlzLmN1cnJlbnRNb250aE5hbWU9dGhpcy5tb250aHNbbi5nZXRNb250aCgpXSkscCh0aGlzLmVsLHRoaXMpLGkmJih5KHtpbnN0YW5jZTp0aGlzfSksdShpKSk7dmFyIG89YT09PW4uZ2V0RnVsbFllYXIoKSYmcj09PW4uZ2V0TW9udGgoKTtyZXR1cm4gb3x8dD91KHRoaXMsbik6b3x8dSh0aGlzLG5ldyBEYXRlKGEsciwxKSksdGhpc31mdW5jdGlvbiBOKGUpe3JldHVybiBJKHRoaXMsZSwhMCl9ZnVuY3Rpb24gXyhlKXtyZXR1cm4gSSh0aGlzLGUpfWZ1bmN0aW9uIEkoZSx0LG4pe3ZhciBhPWUuZGF0ZVNlbGVjdGVkLHI9ZS5maXJzdCxpPWUuc2libGluZyxvPWUubWluRGF0ZSxzPWUubWF4RGF0ZSxsPWcodCksZD1uPyJNaW4iOiJNYXgiO2Z1bmN0aW9uIGMoKXtyZXR1cm4ib3JpZ2luYWwiK2QrIkRhdGUifWZ1bmN0aW9uIGgoKXtyZXR1cm4gZC50b0xvd2VyQ2FzZSgpKyJEYXRlIn1mdW5jdGlvbiBmKCl7cmV0dXJuInNldCIrZH1mdW5jdGlvbiB2KCl7dGhyb3cgbmV3IEVycm9yKCJPdXQtb2YtcmFuZ2UgZGF0ZSBwYXNzZWQgdG8gIitmKCkpfWlmKG51bGw9PXQpZVtjKCldPXZvaWQgMCxpPyhpW2MoKV09dm9pZCAwLG4/KHImJiFhfHwhciYmIWkuZGF0ZVNlbGVjdGVkKSYmKGUubWluRGF0ZT12b2lkIDAsaS5taW5EYXRlPXZvaWQgMCk6KHImJiFpLmRhdGVTZWxlY3RlZHx8IXImJiFhKSYmKGUubWF4RGF0ZT12b2lkIDAsaS5tYXhEYXRlPXZvaWQgMCkpOmVbaCgpXT12b2lkIDA7ZWxzZXtpZighYih0KSl0aHJvdyBuZXcgRXJyb3IoIkludmFsaWQgZGF0ZSBwYXNzZWQgdG8gIitmKCkpO2k/KChyJiZuJiZsPihhfHxzKXx8ciYmIW4mJmw8KGkuZGF0ZVNlbGVjdGVkfHxvKXx8IXImJm4mJmw+KGkuZGF0ZVNlbGVjdGVkfHxzKXx8IXImJiFuJiZsPChhfHxvKSkmJnYoKSxlW2MoKV09bCxpW2MoKV09bCwobiYmKHImJiFhfHwhciYmIWkuZGF0ZVNlbGVjdGVkKXx8IW4mJihyJiYhaS5kYXRlU2VsZWN0ZWR8fCFyJiYhYSkpJiYoZVtoKCldPWwsaVtoKCldPWwpKTooKG4mJmw+KGF8fHMpfHwhbiYmbDwoYXx8bykpJiZ2KCksZVtoKCldPWwpfXJldHVybiBpJiZ1KGkpLHUoZSksZX1mdW5jdGlvbiBBKCl7dmFyIGU9dGhpcy5maXJzdD90aGlzOnRoaXMuc2libGluZyx0PWUuc2libGluZztyZXR1cm57c3RhcnQ6ZS5kYXRlU2VsZWN0ZWQsZW5kOnQuZGF0ZVNlbGVjdGVkfX1mdW5jdGlvbiBSKCl7dmFyIGU9dGhpcy5zaGFkb3dEb20sdD10aGlzLnBvc2l0aW9uZWRFbCxuPXRoaXMuY2FsZW5kYXJDb250YWluZXIscj10aGlzLnNpYmxpbmcsaT10aGlzO3RoaXMuaW5saW5lUG9zaXRpb24mJihhLnNvbWUoKGZ1bmN0aW9uKGUpe3JldHVybiBlIT09aSYmZS5wb3NpdGlvbmVkRWw9PT10fSkpfHx0LnN0eWxlLnNldFByb3BlcnR5KCJwb3NpdGlvbiIsbnVsbCkpO24ucmVtb3ZlKCksYT1hLmZpbHRlcigoZnVuY3Rpb24oZSl7cmV0dXJuIGUhPT1pfSkpLHImJmRlbGV0ZSByLnNpYmxpbmcsYS5sZW5ndGh8fGooZG9jdW1lbnQsTCk7dmFyIG89YS5zb21lKChmdW5jdGlvbih0KXtyZXR1cm4gdC5zaGFkb3dEb209PT1lfSkpO2Zvcih2YXIgcyBpbiBlJiYhbyYmaihlLFkpLHRoaXMpZGVsZXRlIHRoaXNbc107YS5sZW5ndGh8fGwuZm9yRWFjaCgoZnVuY3Rpb24oZSl7ZG9jdW1lbnQucmVtb3ZlRXZlbnRMaXN0ZW5lcihlLEwpfSkpfWZ1bmN0aW9uIEYoZSx0KXt2YXIgbj1uZXcgRGF0ZShlKTtpZighYihuKSl0aHJvdyBuZXcgRXJyb3IoIkludmFsaWQgZGF0ZSBwYXNzZWQgdG8gYG5hdmlnYXRlYCIpO3RoaXMuY3VycmVudFllYXI9bi5nZXRGdWxsWWVhcigpLHRoaXMuY3VycmVudE1vbnRoPW4uZ2V0TW9udGgoKSx1KHRoaXMpLHQmJnRoaXMub25Nb250aENoYW5nZSh0aGlzKX1mdW5jdGlvbiBCKCl7dmFyIGU9IXRoaXMuY2FsZW5kYXJDb250YWluZXIuY2xhc3NMaXN0LmNvbnRhaW5zKCJxcy1oaWRkZW4iKSx0PSF0aGlzLmNhbGVuZGFyQ29udGFpbmVyLnF1ZXJ5U2VsZWN0b3IoIi5xcy1vdmVybGF5IikuY2xhc3NMaXN0LmNvbnRhaW5zKCJxcy1oaWRkZW4iKTtlJiZNKHQsdGhpcyl9dC5kZWZhdWx0PWZ1bmN0aW9uKGUsdCl7dmFyIG49ZnVuY3Rpb24oZSx0KXt2YXIgbixsLGQ9ZnVuY3Rpb24oZSl7dmFyIHQ9YyhlKTt0LmV2ZW50cyYmKHQuZXZlbnRzPXQuZXZlbnRzLnJlZHVjZSgoZnVuY3Rpb24oZSx0KXtpZighYih0KSl0aHJvdyBuZXcgRXJyb3IoJyJvcHRpb25zLmV2ZW50cyIgbXVzdCBvbmx5IGNvbnRhaW4gdmFsaWQgSmF2YVNjcmlwdCBEYXRlIG9iamVjdHMuJyk7cmV0dXJuIGVbK2codCldPSEwLGV9KSx7fSkpO1sic3RhcnREYXRlIiwiZGF0ZVNlbGVjdGVkIiwibWluRGF0ZSIsIm1heERhdGUiXS5mb3JFYWNoKChmdW5jdGlvbihlKXt2YXIgbj10W2VdO2lmKG4mJiFiKG4pKXRocm93IG5ldyBFcnJvcignIm9wdGlvbnMuJytlKyciIG5lZWRzIHRvIGJlIGEgdmFsaWQgSmF2YVNjcmlwdCBEYXRlIG9iamVjdC4nKTt0W2VdPWcobil9KSk7dmFyIG49dC5wb3NpdGlvbixpPXQubWF4RGF0ZSxsPXQubWluRGF0ZSxkPXQuZGF0ZVNlbGVjdGVkLHU9dC5vdmVybGF5UGxhY2Vob2xkZXIsaD10Lm92ZXJsYXlCdXR0b24sZj10LnN0YXJ0RGF5LHY9dC5pZDtpZih0LnN0YXJ0RGF0ZT1nKHQuc3RhcnREYXRlfHxkfHxuZXcgRGF0ZSksdC5kaXNhYmxlZERhdGVzPSh0LmRpc2FibGVkRGF0ZXN8fFtdKS5yZWR1Y2UoKGZ1bmN0aW9uKGUsdCl7dmFyIG49K2codCk7aWYoIWIodCkpdGhyb3cgbmV3IEVycm9yKCdZb3Ugc3VwcGxpZWQgYW4gaW52YWxpZCBkYXRlIHRvICJvcHRpb25zLmRpc2FibGVkRGF0ZXMiLicpO2lmKG49PT0rZyhkKSl0aHJvdyBuZXcgRXJyb3IoJyJkaXNhYmxlZERhdGVzIiBjYW5ub3QgY29udGFpbiB0aGUgc2FtZSBkYXRlIGFzICJkYXRlU2VsZWN0ZWQiLicpO3JldHVybiBlW25dPTEsZX0pLHt9KSx0Lmhhc093blByb3BlcnR5KCJpZCIpJiZudWxsPT12KXRocm93IG5ldyBFcnJvcigiYGlkYCBjYW5ub3QgYmUgYG51bGxgIG9yIGB1bmRlZmluZWRgIik7aWYobnVsbCE9dil7dmFyIG09YS5maWx0ZXIoKGZ1bmN0aW9uKGUpe3JldHVybiBlLmlkPT09dn0pKTtpZihtLmxlbmd0aD4xKXRocm93IG5ldyBFcnJvcigiT25seSB0d28gZGF0ZXBpY2tlcnMgY2FuIHNoYXJlIGFuIGlkLiIpO20ubGVuZ3RoPyh0LnNlY29uZD0hMCx0LnNpYmxpbmc9bVswXSk6dC5maXJzdD0hMH12YXIgeT1bInRyIiwidGwiLCJiciIsImJsIiwiYyJdLnNvbWUoKGZ1bmN0aW9uKGUpe3JldHVybiBuPT09ZX0pKTtpZihuJiYheSl0aHJvdyBuZXcgRXJyb3IoJyJvcHRpb25zLnBvc2l0aW9uIiBtdXN0IGJlIG9uZSBvZiB0aGUgZm9sbG93aW5nOiB0bCwgdHIsIGJsLCBiciwgb3IgYy4nKTtmdW5jdGlvbiBwKGUpe3Rocm93IG5ldyBFcnJvcignImRhdGVTZWxlY3RlZCIgaW4gb3B0aW9ucyBpcyAnKyhlPyJsZXNzIjoiZ3JlYXRlciIpKycgdGhhbiAiJysoZXx8Im1heCIpKydEYXRlIi4nKX1pZih0LnBvc2l0aW9uPWZ1bmN0aW9uKGUpe3ZhciB0PWVbMF0sbj1lWzFdLGE9e307YVtvW3RdXT0xLG4mJihhW29bbl1dPTEpO3JldHVybiBhfShufHwiYmwiKSxpPGwpdGhyb3cgbmV3IEVycm9yKCcibWF4RGF0ZSIgaW4gb3B0aW9ucyBpcyBsZXNzIHRoYW4gIm1pbkRhdGUiLicpO2QmJihsPmQmJnAoIm1pbiIpLGk8ZCYmcCgpKTtpZihbIm9uU2VsZWN0Iiwib25TaG93Iiwib25IaWRlIiwib25Nb250aENoYW5nZSIsImZvcm1hdHRlciIsImRpc2FibGVyIl0uZm9yRWFjaCgoZnVuY3Rpb24oZSl7ImZ1bmN0aW9uIiE9dHlwZW9mIHRbZV0mJih0W2VdPXMpfSkpLFsiY3VzdG9tRGF5cyIsImN1c3RvbU1vbnRocyIsImN1c3RvbU92ZXJsYXlNb250aHMiXS5mb3JFYWNoKChmdW5jdGlvbihlLG4pe3ZhciBhPXRbZV0scj1uPzEyOjc7aWYoYSl7aWYoIUFycmF5LmlzQXJyYXkoYSl8fGEubGVuZ3RoIT09cnx8YS5zb21lKChmdW5jdGlvbihlKXtyZXR1cm4ic3RyaW5nIiE9dHlwZW9mIGV9KSkpdGhyb3cgbmV3IEVycm9yKCciJytlKyciIG11c3QgYmUgYW4gYXJyYXkgd2l0aCAnK3IrIiBzdHJpbmdzLiIpO3Rbbj9uPDI/Im1vbnRocyI6Im92ZXJsYXlNb250aHMiOiJkYXlzIl09YX19KSksZiYmZj4wJiZmPDcpe3ZhciB3PSh0LmN1c3RvbURheXN8fHIpLnNsaWNlKCksRD13LnNwbGljZSgwLGYpO3QuY3VzdG9tRGF5cz13LmNvbmNhdChEKSx0LnN0YXJ0RGF5PStmLHQud2Vla2VuZEluZGljZXM9W3cubGVuZ3RoLTEsdy5sZW5ndGhdfWVsc2UgdC5zdGFydERheT0wLHQud2Vla2VuZEluZGljZXM9WzYsMF07InN0cmluZyIhPXR5cGVvZiB1JiZkZWxldGUgdC5vdmVybGF5UGxhY2Vob2xkZXI7InN0cmluZyIhPXR5cGVvZiBoJiZkZWxldGUgdC5vdmVybGF5QnV0dG9uO3ZhciBxPXQuZGVmYXVsdFZpZXc7aWYocSYmImNhbGVuZGFyIiE9PXEmJiJvdmVybGF5IiE9PXEpdGhyb3cgbmV3IEVycm9yKCdvcHRpb25zLmRlZmF1bHRWaWV3IG11c3QgZWl0aGVyIGJlICJjYWxlbmRhciIgb3IgIm92ZXJsYXkiLicpO3JldHVybiB0LmRlZmF1bHRWaWV3PXF8fCJjYWxlbmRhciIsdH0odHx8e3N0YXJ0RGF0ZTpnKG5ldyBEYXRlKSxwb3NpdGlvbjoiYmwiLGRlZmF1bHRWaWV3OiJjYWxlbmRhciJ9KSx1PWU7aWYoInN0cmluZyI9PXR5cGVvZiB1KXU9IiMiPT09dVswXT9kb2N1bWVudC5nZXRFbGVtZW50QnlJZCh1LnNsaWNlKDEpKTpkb2N1bWVudC5xdWVyeVNlbGVjdG9yKHUpO2Vsc2V7aWYoIltvYmplY3QgU2hhZG93Um9vdF0iPT09eCh1KSl0aHJvdyBuZXcgRXJyb3IoIlVzaW5nIGEgc2hhZG93IERPTSBhcyB5b3VyIHNlbGVjdG9yIGlzIG5vdCBzdXBwb3J0ZWQuIik7Zm9yKHZhciBoLGY9dS5wYXJlbnROb2RlOyFoOyl7dmFyIHY9eChmKTsiW29iamVjdCBIVE1MRG9jdW1lbnRdIj09PXY/aD0hMDoiW29iamVjdCBTaGFkb3dSb290XSI9PT12PyhoPSEwLG49ZixsPWYuaG9zdCk6Zj1mLnBhcmVudE5vZGV9fWlmKCF1KXRocm93IG5ldyBFcnJvcigiTm8gc2VsZWN0b3IgLyBlbGVtZW50IGZvdW5kLiIpO2lmKGEuc29tZSgoZnVuY3Rpb24oZSl7cmV0dXJuIGUuZWw9PT11fSkpKXRocm93IG5ldyBFcnJvcigiQSBkYXRlcGlja2VyIGFscmVhZHkgZXhpc3RzIG9uIHRoYXQgZWxlbWVudC4iKTt2YXIgbT11PT09ZG9jdW1lbnQuYm9keSx5PW4/dS5wYXJlbnRFbGVtZW50fHxuOm0/ZG9jdW1lbnQuYm9keTp1LnBhcmVudEVsZW1lbnQsdz1uP3UucGFyZW50RWxlbWVudHx8bDp5LEQ9ZG9jdW1lbnQuY3JlYXRlRWxlbWVudCgiZGl2IikscT1kb2N1bWVudC5jcmVhdGVFbGVtZW50KCJkaXYiKTtELmNsYXNzTmFtZT0icXMtZGF0ZXBpY2tlci1jb250YWluZXIgcXMtaGlkZGVuIixxLmNsYXNzTmFtZT0icXMtZGF0ZXBpY2tlciI7dmFyIE09e3NoYWRvd0RvbTpuLGN1c3RvbUVsZW1lbnQ6bCxwb3NpdGlvbmVkRWw6dyxlbDp1LHBhcmVudDp5LG5vbklucHV0OiJJTlBVVCIhPT11Lm5vZGVOYW1lLG5vUG9zaXRpb246bSxwb3NpdGlvbjohbSYmZC5wb3NpdGlvbixzdGFydERhdGU6ZC5zdGFydERhdGUsZGF0ZVNlbGVjdGVkOmQuZGF0ZVNlbGVjdGVkLGRpc2FibGVkRGF0ZXM6ZC5kaXNhYmxlZERhdGVzLG1pbkRhdGU6ZC5taW5EYXRlLG1heERhdGU6ZC5tYXhEYXRlLG5vV2Vla2VuZHM6ISFkLm5vV2Vla2VuZHMsd2Vla2VuZEluZGljZXM6ZC53ZWVrZW5kSW5kaWNlcyxjYWxlbmRhckNvbnRhaW5lcjpELGNhbGVuZGFyOnEsY3VycmVudE1vbnRoOihkLnN0YXJ0RGF0ZXx8ZC5kYXRlU2VsZWN0ZWQpLmdldE1vbnRoKCksY3VycmVudE1vbnRoTmFtZTooZC5tb250aHN8fGkpWyhkLnN0YXJ0RGF0ZXx8ZC5kYXRlU2VsZWN0ZWQpLmdldE1vbnRoKCldLGN1cnJlbnRZZWFyOihkLnN0YXJ0RGF0ZXx8ZC5kYXRlU2VsZWN0ZWQpLmdldEZ1bGxZZWFyKCksZXZlbnRzOmQuZXZlbnRzfHx7fSxkZWZhdWx0VmlldzpkLmRlZmF1bHRWaWV3LHNldERhdGU6TyxyZW1vdmU6UixzZXRNaW46TixzZXRNYXg6XyxzaG93OlAsaGlkZTprLG5hdmlnYXRlOkYsdG9nZ2xlT3ZlcmxheTpCLG9uU2VsZWN0OmQub25TZWxlY3Qsb25TaG93OmQub25TaG93LG9uSGlkZTpkLm9uSGlkZSxvbk1vbnRoQ2hhbmdlOmQub25Nb250aENoYW5nZSxmb3JtYXR0ZXI6ZC5mb3JtYXR0ZXIsZGlzYWJsZXI6ZC5kaXNhYmxlcixtb250aHM6ZC5tb250aHN8fGksZGF5czpkLmN1c3RvbURheXN8fHIsc3RhcnREYXk6ZC5zdGFydERheSxvdmVybGF5TW9udGhzOmQub3ZlcmxheU1vbnRoc3x8KGQubW9udGhzfHxpKS5tYXAoKGZ1bmN0aW9uKGUpe3JldHVybiBlLnNsaWNlKDAsMyl9KSksb3ZlcmxheVBsYWNlaG9sZGVyOmQub3ZlcmxheVBsYWNlaG9sZGVyfHwiNC1kaWdpdCB5ZWFyIixvdmVybGF5QnV0dG9uOmQub3ZlcmxheUJ1dHRvbnx8IlN1Ym1pdCIsZGlzYWJsZVllYXJPdmVybGF5OiEhZC5kaXNhYmxlWWVhck92ZXJsYXksZGlzYWJsZU1vYmlsZTohIWQuZGlzYWJsZU1vYmlsZSxpc01vYmlsZToib250b3VjaHN0YXJ0ImluIHdpbmRvdyxhbHdheXNTaG93OiEhZC5hbHdheXNTaG93LGlkOmQuaWQsc2hvd0FsbERhdGVzOiEhZC5zaG93QWxsRGF0ZXMscmVzcGVjdERpc2FibGVkUmVhZE9ubHk6ISFkLnJlc3BlY3REaXNhYmxlZFJlYWRPbmx5LGZpcnN0OmQuZmlyc3Qsc2Vjb25kOmQuc2Vjb25kfTtpZihkLnNpYmxpbmcpe3ZhciBFPWQuc2libGluZyxDPU0sTD1FLm1pbkRhdGV8fEMubWluRGF0ZSxZPUUubWF4RGF0ZXx8Qy5tYXhEYXRlO0Muc2libGluZz1FLEUuc2libGluZz1DLEUubWluRGF0ZT1MLEUubWF4RGF0ZT1ZLEMubWluRGF0ZT1MLEMubWF4RGF0ZT1ZLEUub3JpZ2luYWxNaW5EYXRlPUwsRS5vcmlnaW5hbE1heERhdGU9WSxDLm9yaWdpbmFsTWluRGF0ZT1MLEMub3JpZ2luYWxNYXhEYXRlPVksRS5nZXRSYW5nZT1BLEMuZ2V0UmFuZ2U9QX1kLmRhdGVTZWxlY3RlZCYmcCh1LE0pO3ZhciBqPWdldENvbXB1dGVkU3R5bGUodykucG9zaXRpb247bXx8aiYmInN0YXRpYyIhPT1qfHwoTS5pbmxpbmVQb3NpdGlvbj0hMCx3LnN0eWxlLnNldFByb3BlcnR5KCJwb3NpdGlvbiIsInJlbGF0aXZlIikpO3ZhciBJPWEuZmlsdGVyKChmdW5jdGlvbihlKXtyZXR1cm4gZS5wb3NpdGlvbmVkRWw9PT1NLnBvc2l0aW9uZWRFbH0pKTtJLnNvbWUoKGZ1bmN0aW9uKGUpe3JldHVybiBlLmlubGluZVBvc2l0aW9ufSkpJiYoTS5pbmxpbmVQb3NpdGlvbj0hMCxJLmZvckVhY2goKGZ1bmN0aW9uKGUpe2UuaW5saW5lUG9zaXRpb249ITB9KSkpO0QuYXBwZW5kQ2hpbGQocSkseS5hcHBlbmRDaGlsZChEKSxNLmFsd2F5c1Nob3cmJlMoTSk7cmV0dXJuIE19KGUsdCk7aWYoYS5sZW5ndGh8fGQoZG9jdW1lbnQpLG4uc2hhZG93RG9tJiYoYS5zb21lKChmdW5jdGlvbihlKXtyZXR1cm4gZS5zaGFkb3dEb209PT1uLnNoYWRvd0RvbX0pKXx8ZChuLnNoYWRvd0RvbSkpLGEucHVzaChuKSxuLnNlY29uZCl7dmFyIGw9bi5zaWJsaW5nO3koe2luc3RhbmNlOm4sZGVzZWxlY3Q6IW4uZGF0ZVNlbGVjdGVkfSkseSh7aW5zdGFuY2U6bCxkZXNlbGVjdDohbC5kYXRlU2VsZWN0ZWR9KSx1KGwpfXJldHVybiB1KG4sbi5zdGFydERhdGV8fG4uZGF0ZVNlbGVjdGVkKSxuLmFsd2F5c1Nob3cmJkQobiksbn19XSkuZGVmYXVsdH0pKTs=`

// embedded https://cdn.jsdelivr.net/npm/js-datepicker@5.18.0/dist/datepicker.min.css
func DatepickerCSS() io.Reader {
	return strings.NewReader(datepickercss)
}

// embedded https://cdn.jsdelivr.net/npm/js-datepicker@5.18.0/dist/datepicker.min.js
func DatepickerJS() io.Reader {
	return base64.NewDecoder(base64.StdEncoding, strings.NewReader(strings.Replace(datepickerjs, "|'|", "`", -1)))
}

func init() {
	gblrs := resources.GLOBALRSNG()
	gblrs.FS().MKDIR("/datepicker/css", "")
	gblrs.FS().MKDIR("/datepicker/js", "")
	gblrs.FS().SET("/datepicker/css/datepicker.css", DatepickerCSS())
	gblrs.FS().SET("/datepicker/js/datepicker.js", DatepickerJS())
}
