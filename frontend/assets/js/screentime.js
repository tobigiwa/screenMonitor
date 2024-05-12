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
    var chartControl = document.querySelectorAll( '.left-controls-button, .month' );
    chartControl.forEach( function ( control )
    {
        control.addEventListener( 'htmx:afterRequest', function ( event )
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
            setCurrentMonth( response.weekStatResponse.month );
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
    let label = `from ${response.weekStatResponse.formattedDay[ 0 ].slice( 5, )} - ${response.weekStatResponse.formattedDay[ 6 ]} ${response.weekStatResponse.month}`;
    currentSaturday = response.weekStatResponse.keys[ 6 ];

    myChart = new Chart( ctx, {
        type: 'bar',
        data: {
            labels: labels,
            datasets: [ {
                label: label,
                data: data,
                borderWidth: 1,
                backgroundColor: [
                    'rgba(255, 99, 132, 0.5)',
                    'rgba(255, 159, 64, 0.5)',
                    'rgba(255, 205, 86, 0.5)',
                    'rgba(75, 192, 192, 0.5)',
                    'rgba(54, 162, 235, 0.5)',
                    'rgba(153, 102, 255, 0.5)',
                    'rgba(201, 203, 207, 0.5)',

                ],
                borderColor: [
                    'rgb(255, 99, 132)',
                    'rgb(255, 159, 64)',
                    'rgb(255, 205, 86)',
                    'rgb(75, 192, 192)',
                    'rgb(54, 162, 235)',
                    'rgb(153, 102, 255)',
                    'rgb(201, 203, 207)',
                ],
            } ]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            scales: {
                y: {
                    beginAtZero: true,
                    title: {
                        display: true,
                        text: 'in Hours',
                    }
                }
            },
            plugins: {
                title: {
                    display: true,
                    text: 'Weekly Screentime'
                },
                tooltip: {
                    usePointStyle: true,
                    callbacks: {
                        labelTextColor: function ( tooltipItem )
                        {
                            return '#ffffff';
                        },
                        labelPointStyle: function ( tooltipItem )
                        {
                            return {
                                pointStyle: 'triangle',
                                rotation: 0
                            };
                        },
                        label: function ( tooltipItem )
                        {
                            // var value = 'Active Uptime: ' + tooltipItem.parsed.y.toFixed( 2 ) + 'Hrs';
                            var value = 'Active Uptime: ' + Number( tooltipItem.parsed.y.toFixed( 2 ) ) + 'Hrs';
                            return value;
                        },
                    }
                }
            }
        }
    } );
}

// function setCurrentMonth ( month )
// {
//     let selectElement = document.querySelector( '.month' );
//     let selectOptions = selectElement.options;

//     // Check if the month is in the dropdown
//     let monthInDropdown = Array.from( selectOptions ).some( option => option.value == month );

//     if ( monthInDropdown ) {
//         selectElement.value = month;
//         selectElement.style.fontSize = "initial";
//         selectElement.style.fontStyle = "initial";
//     } else {
//         selectElement.value = "";
//         selectElement.style.fontSize = "smaller";
//         selectElement.style.fontStyle = "italic";
//     }
// }


function setCurrentMonth ( month )
{
    let selectElement = document.querySelector( '.month' );
    let selectOptions = selectElement.options;

    for ( let option of selectOptions ) {
        option.style.fontSize = "initial";
        option.style.fontStyle = "normal";
    }

    let monthInDropdown = Array.from( selectOptions ).some( option => option.value == month );

    if ( monthInDropdown ) {
        selectElement.value = month;
    } else {
        let placeholderOption = selectElement.querySelector( '#placeholder' );
        placeholderOption.style.fontSize = "smaller";
        placeholderOption.style.fontStyle = "italic";
        selectElement.value = "";
    }
}

function arrowButton ()
{

    let backwardButton = document.querySelector( '.backward-arrow' );
    let endPointA = '/weekStat?week=backward-arrow&saturday=' + currentSaturday;

    backwardButton.setAttribute( 'hx-get', endPointA );
    htmx.process( backwardButton );

    let forwardButton = document.querySelector( '.forward-arrow' );
    let endPointB = '/weekStat?week=forward-arrow&saturday=' + currentSaturday;

    forwardButton.setAttribute( 'hx-get', endPointB );
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
