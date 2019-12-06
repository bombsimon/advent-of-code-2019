#!/usr/bin/env perl

use strict;
use warnings;
use utf8;
use feature qw( say );

use Carp;

=head1 DESCRIPTION

This is the solution for part 1 and 2 of day 6.

=head2 main

Man method to run the program

=cut

sub main {
    my $file = $ARGV[0] || croak 'No input file';

    open( my $fh, '<', $file ) || croak 'Could not open file: ' . $!;
    my @lines = <$fh>;
    close( $fh );

    my %orbit   = ();
    my %backref = ();

    # Map all orbits and keep a backref to where the current orbit came from.
    foreach my $relation ( @lines ) {
        chomp $relation;

        my ( $lhs, $rhs ) = split( /\)/x, $relation );
        $orbit{$lhs}{$rhs} = 1;
        $backref{$rhs} = $lhs;
    }

    my %distances = ();
    my $total     = 0;

    # Iterate through all orbits and traverse them back to 'COM' to get a hash
    # back with each planet as key and the distance to 'COM' as value.
    find_distances_from( 'COM', 0, \%orbit, \%distances );

    # Sum all the distances.
    map { $total += $distances{$_} } keys %distances;

    # Part one result.
    say "Total orbits: $total";

    # Store all backreferences from 'YOU' and 'SAN' in one hash each with their
    # backrefs as key.
    my $backref_src = get_backrefs_from( 'YOU', \%backref, {} );
    my $backref_dst = get_backrefs_from( 'SAN', \%backref, {} );

    # Since backref only contains one planet we know that the number of planets
    # seen is the number of steps back to 'COM'.
    my $src_distances_from_com = scalar keys %$backref_src;
    my $dst_distances_from_com = scalar keys %$backref_dst;

    # Now we just need to find all the common planets (where the orbital courses
    # meet)
    my @common_keys = ();
    foreach my $key ( keys %$backref_src ) {
        push @common_keys, $key if exists $backref_dst->{$key};
    }

    # Check which common planet is the furthest out from 'COM' (and thus closest
    # to us).
    my $nearest_shared_distance = 0;
    foreach my $key ( @common_keys ) {
        $nearest_shared_distance = $distances{$key}
          if $distances{$key} > $nearest_shared_distance || $nearest_shared_distance == 0;
    }

    # Take the delta of the nearest sahred and our planets and add them together
    # to get the distance between 'YOU' and 'SAN'.
    my $delta_src = $src_distances_from_com - $nearest_shared_distance;
    my $delta_dst = $dst_distances_from_com - $nearest_shared_distance;
    my $delta     = $delta_src + $delta_dst - 2;                       # Remove the planets themself.

    # Part two result.
    say "Minimum orbits to move from 'YOU' to 'SAN': $delta";

    return 1;
}

=head2 get_backrefs_from

Recursive function to get a flat hash of all the planets behind $planet.

=cut

sub get_backrefs_from {
    my ( $planet, $backref, $list ) = @_;

    my $backrefed_planet = $backref->{$planet};

    return $list unless $backrefed_planet;

    $list->{$backrefed_planet} = 1;

    {
        no warnings 'recursion';
        get_backrefs_from( $backrefed_planet, $backref, $list );
    }

    return $list;
}

=head2 find_distances_from

Recursive function to get a flat hash of the distance to inputted planet for
each child.

=cut

sub find_distances_from {
    my ( $planet, $distance, $orbit, $distances ) = @_;

    $distances->{$planet} = $distance;

    my $next_distnace = ++$distance;

    foreach my $child ( keys %{ $orbit->{$planet} } ) {
        {
            no warnings 'recursion';
            find_distances_from( $child, $next_distnace, $orbit, $distances );
        };
    }

    return $distances;
}

main();
