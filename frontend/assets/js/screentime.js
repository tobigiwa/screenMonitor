
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
    const ctx = document.getElementById( 'chart' );
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