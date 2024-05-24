import htmx from 'htmx.org';


let myChart;
let currentSaturday;

document.addEventListener( 'DOMContentLoaded', function ()
{

    // onFullPageReload();
    // leftcontrolchartButtons();

} );




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
            htmx.trigger( "#highlight", 'highlight', );
        } );
    } );
}




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

