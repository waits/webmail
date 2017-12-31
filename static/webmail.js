document.addEventListener('DOMContentLoaded', function() {
    const navTop = document.getElementById('navbar').getBoundingClientRect().y;
    let floating = false;
    document.addEventListener('scroll', function() {
        if (!floating && window.pageYOffset >= navTop) {
            document.body.classList.add('has-docked-nav');
            floating = true;
        } else if (floating && window.pageYOffset < navTop) {
            document.body.classList.remove('has-docked-nav');
            floating = false;
        }
    });
});
