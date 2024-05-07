import lottie from 'lottie-web';
import 'htmx.org';
// import * as echarts from 'echarts';
import Chart from 'chart.js/auto';

document.addEventListener( 'DOMContentLoaded', function ()
{
    lottie.loadAnimation( {
        container: document.getElementById( 'lottie-animation' ),
        renderer: 'svg',
        loop: true,
        autoplay: true,
        path: "assets/animation/Animation - 1712666371830.json"
    } );
} );

// document.addEventListener( 'DOMContentLoaded', function ()
// {
//     var chart = document.getElementById( 'echart' );
//     var myChart = echarts.init( chart );
//     window.onresize = function ()
//     {
//         myChart.resize();
//     };
// } );


document.addEventListener( 'keydown', function ( e )
{
    const focusedElement = document.activeElement;
    if ( !focusedElement.classList.contains( 'links' ) ) {
        return;
    }

    let toFocus = null;
    switch ( e.key ) {
        case 'ArrowDown':
            toFocus = focusedElement.parentElement.nextElementSibling;
            if ( toFocus ) toFocus = toFocus.querySelector( '.links' );
            break;
        case 'ArrowUp':
            toFocus = focusedElement.parentElement.previousElementSibling;
            if ( toFocus ) toFocus = toFocus.querySelector( '.links' );
            break;
    }

    if ( toFocus ) {
        toFocus.focus();
        e.preventDefault();
    }
} );

document.addEventListener( 'DOMContentLoaded', function ()
{
    const ctx = document.getElementById( 'echart' );
    new Chart( ctx, {
        type: 'bar',
        data: {
            labels: [ 'Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday' ],
            datasets: [ {
                label: 'Days of the week',
                data: [ 12, 19, 3, 5, 2, 3, 11 ],
                borderWidth: 1,
                backgroundColor: 'rgba(255, 99, 132, 0.2)',
                borderColor: 'rgba(255, 99, 132, 1)',
            } ]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            scales: {
                y: {
                    beginAtZero: true
                }
            }
        }
    } );
} );


document.addEventListener( 'DOMContentLoaded', function ()
{
    console.log( 'htmx:afterRequest event fired' );
    var chartControlButtons = document.querySelectorAll( '.left-controls-button' );

    chartControlButtons.forEach( function ( button )
    {
        button.addEventListener( 'htmx:afterRequest', function ( event )
        {
            if ( !event.detail.successful ) {
                console.log( "request not successful" );
                return;
            }

            console.log( 'AJAX request has finished for button: ', button );
            console.log( 'AJAX request status: ', event.detail.successful );
            console.log( 'XMLHttpRequest: ', event.detail.xhr );
            console.log( 'XMLHttpRequest response: ', event.detail.xhr.response );
        } );
    } );


} );
