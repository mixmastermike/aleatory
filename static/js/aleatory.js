
$(function() {

  // Helper function to daw on the screen and play the note
  play = function( type, size ){
    if( type == 0 ) {
      playSound( type, size );
      cntx.spawn( type, size );
    }
  };

  // ----------------------------------------
  // Particle definition
  // ----------------------------------------

  var Particle = function( x, y, radius ) {
  	this.init( x, y, radius );
  }

  Particle.prototype = {

  	init: function( x, y, radius ) {
  		this.alive = true;
  		this.radius = radius || 10;
  		this.wander = 0.15;
  		this.theta = random( TWO_PI );
  		this.drag = 0.92;
  		this.color = '#fff';
  		this.x = x || 0.0;
  		this.y = y || 0.0;
  		this.vx = 0.0;
  		this.vy = 0.0;
  	},

  	move: function() {
  		this.x += this.vx;
  		this.y += this.vy;
  		this.vx *= this.drag;
  		this.vy *= this.drag;
  		this.theta += random( -0.5, 0.5 ) * this.wander;
  		this.vx += sin( this.theta ) * 0.1;
  		this.vy += cos( this.theta ) * 0.1;
  		this.radius *= 0.96;
  		this.alive = this.radius > 0.5;
  	},

  	draw: function( ctx ) {
  		ctx.beginPath();
  		ctx.arc( this.x, this.y, this.radius, 0, TWO_PI );
  		ctx.fillStyle = this.color;
  		ctx.fill();
  	}
  };

  // ----------------------------------------
  // The particle generator
  // ----------------------------------------

  var MAX_PARTICLES = 256;
  //var COLOURS = [ '#69D2E7', '#A7DBD8', '#E0E4CC', '#F38630', '#FA6900', '#FF4E50', '#F9D423' ];
  var COLOURS = [ '#F6F6F6', '#A2CAD3', '#0F3BA9', '#FF6004', '#C96857', '#FFAEBC', '#C3B382', '#FF3401', '#44443A' ];
  var GREY_SCALE = [ '#AAAAAA', '#BBBBBB', '#CCCCCC', '#DDDDDD', '#EEEEEE' ];

  var particles = [];
  var pool = [];

  // Create the canvas
  var cntx = Sketch.create({
  	container: document.getElementById( 'container' )
  });

  cntx.setup = function() {};

  cntx.spawn = function( type, size ) {

    if ( particles.length >= MAX_PARTICLES )
  		pool.push( particles.shift() );

    // Spawn the particle roughly in the middle of the canvas
    var w = cntx.width/3;
    var h = cntx.height/3;
    var x = random( w, w+w );
    var y = random( h, h+h );

    // Create a new particle, or re-use an old one
  	particle = pool.length ? pool.pop() : new Particle();
  	particle.init( x, y, size );
  	particle.wander = random( 0.5, 2.0 );

    switch( type ) {
      case '0':
        particle.color = random( COLOURS );
        break;
      case '1':
        particle.color = random( GREY_SCALE );
        break;
    }

  	particle.drag = random( 0.9, 0.99 );

    theta = random( TWO_PI );
  	force = random( 2, 8 );
  	particle.vx = sin( theta ) * force;
  	particle.vy = cos( theta ) * force;
    // Add the particle into the generator
  	particles.push( particle );
  }

  // Update the position of particles.  If needed, remove dead ones.
  cntx.update = function() {
  	var i, particle;
  	for ( i = particles.length - 1; i >= 0; i-- ) {
  		particle = particles[i];
      // Remove dead particles
  		if ( particle.alive ) particle.move();
  		else pool.push( particles.splice( i, 1 )[0] );
  	}
  };

  // (re)Draw all the particles on the canvas
  cntx.draw = function() {
  	cntx.globalCompositeOperation  = 'lighter';
  	for ( var i = particles.length - 1; i >= 0; i-- ) {
  		particles[i].draw( cntx );
  	}
  };

  // ----------------------------------------
  // Sounds
  // ----------------------------------------

  var sounds = [];
  var ambience = [];
  ambience[0] = new Howl({
      urls : ['sounds/birds.mp3'],
      volume : 0.05,
      loop: true,
  });
  ambience[1] = new Howl({
      urls : ['sounds/ambient.mp3'],
      volume : 0.1,
      loop: true,
  });
  // Auto-start the ambience
  ambience[0].play();
  ambience[1].play();

  // Load the musical sounds
  for (var i = 10; i <= 24; i++) {
    sounds.push(new Howl({
        urls : ['sounds/harp-'+i+'.wav'],
        volume : 0.2
    }));
  }

  // Play a sound based on the type and size given
  function playSound(type, size) {
    var pitch = Math.round(size / (255 / sounds.length));
    if( sounds[pitch] )
      sounds[pitch].play();
  }


});
