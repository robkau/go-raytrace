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

//There's just one more shape to which you need to add support for divide(): the CSG shape. In this case the divide() function does nothing to the shape itself, but, propagates the call to the shape's children.
//
//The following test demonstrates this with a CSG shape consisting of two groups, each containing two spheres. After calling divide() on the CSG shape, the left and right children will have been split into two subgroups.
//
//csg.feature
//Scenario: Subdividing a CSG shape subdivides its children
//  Given s1 ← sphere() with:
//      | transform | translation(-1.5, 0, 0) |
//    And s2 ← sphere() with:
//      | transform | translation(1.5, 0, 0) |
//    And left ← group() of [s1, s2]
//    And s3 ← sphere() with:
//      | transform | translation(0, 0, -1.5) |
//    And s4 ← sphere() with:
//      | transform | translation(0, 0, 1.5) |
//    And right ← group() of [s3, s4]
//    And shape ← csg("difference", left, right)
//  When divide(shape, 1)
//  Then left[0] is a group of [s1]
//    And left[1] is a group of [s2]
//    And right[0] is a group of [s3]
//    And right[1] is a group of [s4]
