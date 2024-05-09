import lottie from 'lottie-web';
import htmx from 'htmx.org';
// import * as echarts from 'echarts';
import Chart from 'chart.js/auto';


// document.addEventListener( 'DOMContentLoaded', function ()
// {
//     var chart = document.getElementById( 'echart' );
//     var myChart = echarts.init( chart );
//     window.onresize = function ()
//     {
//         myChart.resize();
//     };
// } );

let myChart;
let currentSaturday;

document.addEventListener( 'DOMContentLoaded', function ()
{

    loadAnimation();
    onFullPageReload();
    leftcontrolchartButtons();

} );

function loadAnimation ()
{
    lottie.loadAnimation( {
        container: document.getElementById( 'lottie-animation' ),
        renderer: 'svg',
        loop: true,
        autoplay: true,
        path: "assets/animation/Animation - 1712666371830.json"
    } );
}

function onFullPageReload ()
{
    window.onload = function ()
    {
        var thisWeekButton = document.querySelector( '#thisWeekButton' );

        if ( thisWeekButton ) {
            thisWeekButton.click();
        }
    };
}

function leftcontrolchartButtons ()
{
    var chartControlButtons = document.querySelectorAll( '.left-controls-button' );
    chartControlButtons.forEach( function ( button )
    {
        button.addEventListener( 'htmx:afterRequest', function ( event )
        {
            if ( !event.detail.successful ) {
                console.log( "request not successful" );
                return;
            }

            console.log( 'XMLHttpRequest response: ', event.detail.xhr.response );
            let response;
            try {
                response = JSON.parse( event.detail.xhr.response );
            } catch ( error ) {
                console.error( 'Invalid response: not a valid JSON string', error );
                return;
            }
            drawWeekStatChar( response );
            arrowButton();
        } );
    } );
}

function drawWeekStatChar ( response )
{
    const ctx = document.getElementById( 'echart' );
    if ( myChart ) {
        myChart.destroy();
    }

    let labels = response.weekStatResponse.formattedDay;
    let data = response.weekStatResponse.values;
    currentSaturday = response.weekStatResponse.keys[ 6 ];

    myChart = new Chart( ctx, {
        type: 'bar',
        data: {
            labels: labels,
            datasets: [ {
                label: 'Days of the week',
                data: data,
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
}

function arrowButton ()
{
    console.log( "in arrowButton", currentSaturday );

    let backwardButton = document.querySelector( '.backward-arrow' );
    let endPointA = '/weekStat?week=backward-arrow&saturday=' + currentSaturday;

    backwardButton.setAttribute( 'hx-get', endPointA );
    backwardButton.setAttribute( 'hx-swap', 'none' );
    htmx.process( backwardButton );

    let forwardButton = document.querySelector( '.forward-arrow' );
    let endPointB = '/weekStat?week=forward-arrow&saturday=' + currentSaturday;

    forwardButton.setAttribute( 'hx-get', endPointB );
    forwardButton.setAttribute( 'hx-swap', 'none' );
    htmx.process( forwardButton );
}

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
