// cube.ch
a cube is a shape
  cube vertices are 8 points
  cube edge length is a float

// marker.ch
a marker is a cube 
  edge length is 0.1
  color is blue

  when active:
    edge length is 0.2
    color is red


origin marker is a marker at origin
  color is yellow

north marker is a marker at [0, 1, 0]
  color is red

  when active:
    color is magenta

// planet.ch
a planet is a sphere
  name is a string

the earth is a planet
  name is "Earth"

// user.ch
a user is a person
  persisted

  name is a string
    access is read only after is set

  password is a string
    access is:
    - write only
    - compare only after is set

  unlock is a string
    access is write only

  authenticated when unlock === password

// profile.ch
a profile is an object
  persisted

  language is a languge
