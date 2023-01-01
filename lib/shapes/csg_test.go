package shapes

// TODO from bounding box chapter
//Scenario: A CSG shape has a bounding box that contains its children
//  Given left ← sphere()
//    And right ← sphere() with:
//      | transform | translation(2, 3, 4) |
//    And shape ← csg("difference", left, right)
//  When box ← bounds_of(shape)
//  Then box.min = point(-1, -1, -1)
//    And box.max = point(3, 4, 5)

// Hint: CSG must find the parent-space bounds of each child object, and then merge them all together into a single bounding box.

//Scenario: Intersecting ray+csg doesn't test children if box is missed
//  Given left ← test_shape()
//    And right ← test_shape()
//    And shape ← csg("difference", left, right)
//    And r ← ray(point(0, 0, -5), vector(0, 1, 0))
//  When xs ← intersect(shape, r)
//  Then left.saved_ray is unset
//    And right.saved_ray is unset
//
//Scenario: Intersecting ray+csg tests children if box is hit
//  Given left ← test_shape()
//    And right ← test_shape()
//    And shape ← csg("difference", left, right)
//    And r ← ray(point(0, 0, -5), vector(0, 0, 1))
//  When xs ← intersect(shape, r)
//  Then left.saved_ray is set
//    And right.saved_ray is set

//In pseudocode, the local_intersect() function for both the Group and CSG shapes should be modified to look something like this:
//
//function local_intersect(shape, ray)
//  if intersects(bounds_of(shape), ray)
//    # perform the usual intersection logic
//    # ...
//  else
//    # nothing intersected
//    return ()
//  end if
//end function
