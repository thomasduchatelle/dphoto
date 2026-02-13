import type {Meta, StoryObj} from '@storybook/nextjs-vite';
import {fn} from 'storybook/test';
import {AlbumCard} from './index';
import {Box, Typography} from '@mui/material';
import {AppBackground} from '@/components/AppLayout/AppBackground';
import {Album, AlbumId} from '@/domains/catalog/language/catalog-state';

const createAlbumId = (owner: string, folderName: string): AlbumId => ({owner, folderName});

const meta = {
    title: 'Components/AlbumCard',
    component: AlbumCard,
    parameters: {
        layout: 'fullscreen',
    },
    tags: ['autodocs'],
} satisfies Meta<typeof AlbumCard>;

export default meta;
type Story = StoryObj<typeof meta>;

// ============================================================================
// Story Wrapper Component
// ============================================================================

const StorySection = ({
    title,
    description,
    children,
}: {
    title: string;
    description?: string;
    children: React.ReactNode;
}) => (
    <Box sx={{mb: 6}}>
        <Typography
            sx={{
                mb: 1,
                fontSize: 16,
                fontWeight: 500,
                color: '#4a9ece',
                textTransform: 'uppercase',
                letterSpacing: '0.1em',
            }}
        >
            {title}
        </Typography>
        {description && (
            <Typography sx={{mb: 3, fontSize: 13, color: 'rgba(255, 255, 255, 0.6)', fontStyle: 'italic'}}>
                {description}
            </Typography>
        )}
        <Box sx={{display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: 3}}>{children}</Box>
    </Box>
);

// ============================================================================
// 1. Default State
// ============================================================================

export const Default: Story = {
    render: () => {
        const album: Album = {
            albumId: createAlbumId('owner-1', 'family-vacation'),
            name: 'Family Vacation',
            start: new Date('2026-07-15'),
            end: new Date('2026-07-22'),
            totalCount: 47,
            temperature: 6.7,
            relativeTemperature: 0.25,
            sharedWith: [],
        };

        return (
            <AppBackground>
                <Box sx={{p: 6, maxWidth: 400}}>
                    <AlbumCard album={album} onClick={fn()} />
                </Box>
            </AppBackground>
        );
    },
};

// ============================================================================
// 2. Shared State
// ============================================================================

export const Shared: Story = {
    render: () => {
        const sharedWithOne: Album = {
            albumId: createAlbumId('owner-1', 'vacation-2026'),
            name: 'Summer Vacation 2026',
            start: new Date('2026-08-01'),
            end: new Date('2026-08-15'),
            totalCount: 156,
            temperature: 10.4,
            relativeTemperature: 0.6,
            sharedWith: [{user: {name: 'Alice Johnson', email: 'alice@example.com', picture: 'https://i.pravatar.cc/150?img=1'}}],
        };

        const sharedWithThree: Album = {
            albumId: createAlbumId('owner-1', 'paris-trip'),
            name: 'Paris Trip',
            start: new Date('2026-03-10'),
            end: new Date('2026-03-17'),
            totalCount: 203,
            temperature: 29.0,
            relativeTemperature: 0.85,
            sharedWith: [
                {user: {name: 'Alice Johnson', email: 'alice@example.com', picture: 'https://i.pravatar.cc/150?img=1'}},
                {user: {name: 'Bob Smith', email: 'bob@example.com', picture: 'https://i.pravatar.cc/150?img=2'}},
                {user: {name: 'Carol White', email: 'carol@example.com'}},
            ],
        };

        return (
            <AppBackground>
                <Box sx={{p: 6}}>
                    <StorySection title="Shared Albums" description="Left: Shared with 1 person · Right: Shared with 3 people">
                        <AlbumCard album={sharedWithOne} onClick={fn()} />
                        <AlbumCard album={sharedWithThree} onClick={fn()} />
                    </StorySection>
                </Box>
            </AppBackground>
        );
    },
};

// ============================================================================
// 3. Temperature Levels
// ============================================================================

export const Temperature: Story = {
    render: () => {
        const albums: Album[] = [
            {
                albumId: createAlbumId('owner-1', 'cold-album'),
                name: 'Weekend Getaway',
                start: new Date('2026-06-01'),
                end: new Date('2026-06-03'),
                totalCount: 12,
                temperature: 4.0,
                relativeTemperature: 0.1,
                sharedWith: [],
            },
            {
                albumId: createAlbumId('owner-1', 'cool-album'),
                name: 'City Tour',
                start: new Date('2026-05-10'),
                end: new Date('2026-05-12'),
                totalCount: 34,
                temperature: 11.3,
                relativeTemperature: 0.35,
                sharedWith: [],
            },
            {
                albumId: createAlbumId('owner-1', 'medium-album'),
                name: 'Beach Holiday',
                start: new Date('2026-07-20'),
                end: new Date('2026-07-27'),
                totalCount: 89,
                temperature: 12.7,
                relativeTemperature: 0.55,
                sharedWith: [],
            },
            {
                albumId: createAlbumId('owner-1', 'warm-album'),
                name: 'Iceland Adventure',
                start: new Date('2026-08-01'),
                end: new Date('2026-08-10'),
                totalCount: 187,
                temperature: 18.7,
                relativeTemperature: 0.82,
                sharedWith: [],
            },
            {
                albumId: createAlbumId('owner-1', 'hot-album'),
                name: 'Christmas 2025',
                start: new Date('2025-12-24'),
                end: new Date('2025-12-26'),
                totalCount: 234,
                temperature: 78.0,
                relativeTemperature: 1.0,
                sharedWith: [],
            },
        ];

        return (
            <AppBackground>
                <Box sx={{p: 6}}>
                    <StorySection
                        title="Temperature Levels"
                        description="Cold (0.1) → Cool (0.35) → Medium (0.55) → Warm (0.82) → Hot (1.0)"
                    >
                        {albums.map((album) => (
                            <AlbumCard key={album.albumId.folderName} album={album} onClick={fn()} />
                        ))}
                    </StorySection>
                </Box>
            </AppBackground>
        );
    },
};

// ============================================================================
// 4. Photo Count Variations (Empty to Few)
// ============================================================================

export const PhotoCount: Story = {
    render: () => {
        const albums: Album[] = [
            {
                albumId: createAlbumId('owner-1', 'empty-album'),
                name: 'Empty Album',
                start: new Date('2026-01-01'),
                end: new Date('2026-01-01'),
                totalCount: 0,
                temperature: 0,
                relativeTemperature: 0.0,
                sharedWith: [],
            },
            {
                albumId: createAlbumId('owner-1', 'one-photo'),
                name: 'Single Photo',
                start: new Date('2026-02-14'),
                end: new Date('2026-02-14'),
                totalCount: 1,
                temperature: 1.0,
                relativeTemperature: 0.05,
                sharedWith: [],
            },
            {
                albumId: createAlbumId('owner-1', 'two-photos'),
                name: 'Two Photos',
                start: new Date('2026-03-08'),
                end: new Date('2026-03-08'),
                totalCount: 2,
                temperature: 2.0,
                relativeTemperature: 0.08,
                sharedWith: [],
            },
            {
                albumId: createAlbumId('owner-1', 'three-photos'),
                name: 'Three Photos',
                start: new Date('2026-04-22'),
                end: new Date('2026-04-22'),
                totalCount: 3,
                temperature: 3.0,
                relativeTemperature: 0.12,
                sharedWith: [],
            },
        ];

        return (
            <AppBackground>
                <Box sx={{p: 6}}>
                    <StorySection title="Photo Count Variations" description="0 photos → 1 photo → 2 photos → 3 photos">
                        {albums.map((album) => (
                            <AlbumCard key={album.albumId.folderName} album={album} onClick={fn()} />
                        ))}
                    </StorySection>
                </Box>
            </AppBackground>
        );
    },
};

// ============================================================================
// 5. Shared By (Album owned by another user)
// ============================================================================

export const SharedBy: Story = {
    render: () => {
        const sharedByOne: Album = {
            albumId: createAlbumId('alice-owner', 'alice-vacation'),
            name: "Alice's Summer Trip",
            start: new Date('2026-06-15'),
            end: new Date('2026-06-22'),
            totalCount: 124,
            temperature: 17.7,
            relativeTemperature: 0.65,
            ownedBy: {
                name: 'Alice Johnson',
                users: [{name: 'Alice Johnson', email: 'alice@example.com', picture: 'https://i.pravatar.cc/150?img=1'}],
            },
            sharedWith: [],
        };

        const sharedByMultiple: Album = {
            albumId: createAlbumId('family-owner', 'family-reunion'),
            name: 'Family Reunion 2026',
            start: new Date('2026-12-20'),
            end: new Date('2026-12-27'),
            totalCount: 287,
            temperature: 35.9,
            relativeTemperature: 0.92,
            ownedBy: {
                name: 'Family Account',
                users: [
                    {name: 'John Doe', email: 'john@example.com', picture: 'https://i.pravatar.cc/150?img=3'},
                    {name: 'Jane Doe', email: 'jane@example.com', picture: 'https://i.pravatar.cc/150?img=4'},
                ],
            },
            sharedWith: [],
        };

        return (
            <AppBackground>
                <Box sx={{p: 6}}>
                    <StorySection
                        title="Shared By (Albums from Other Users)"
                        description="Left: Shared by one user · Right: Shared by multiple users (family account)"
                    >
                        <AlbumCard album={sharedByOne} onClick={fn()} />
                        <AlbumCard album={sharedByMultiple} onClick={fn()} />
                    </StorySection>
                </Box>
            </AppBackground>
        );
    },
};

// ============================================================================
// 6. Combined States
// ============================================================================

export const CombinedStates: Story = {
    render: () => {
        const albums: Album[] = [
            // Default (not shared, low temp)
            {
                albumId: createAlbumId('owner-1', 'default'),
                name: 'Morning Walk',
                start: new Date('2026-05-01'),
                end: new Date('2026-05-01'),
                totalCount: 8,
                temperature: 8.0,
                relativeTemperature: 0.2,
                sharedWith: [],
            },
            // Shared + high temp
            {
                albumId: createAlbumId('owner-1', 'shared-hot'),
                name: 'New York Trip',
                start: new Date('2026-09-10'),
                end: new Date('2026-09-14'),
                totalCount: 312,
                temperature: 62.4,
                relativeTemperature: 0.95,
                sharedWith: [
                    {user: {name: 'Bob Smith', email: 'bob@example.com'}},
                    {user: {name: 'Carol White', email: 'carol@example.com'}},
                ],
            },
            // Shared by + medium temp
            {
                albumId: createAlbumId('bob-owner', 'bob-album'),
                name: "Bob's Birthday Party",
                start: new Date('2026-07-15'),
                end: new Date('2026-07-15'),
                totalCount: 45,
                temperature: 45.0,
                relativeTemperature: 0.58,
                ownedBy: {
                    name: 'Bob Smith',
                    users: [{name: 'Bob Smith', email: 'bob@example.com', picture: 'https://i.pravatar.cc/150?img=2'}],
                },
                sharedWith: [],
            },
            // Long name + shared
            {
                albumId: createAlbumId('owner-1', 'long-name'),
                name: 'Very Long Album Name That Should Truncate With Ellipsis',
                start: new Date('2026-04-01'),
                end: new Date('2026-04-10'),
                totalCount: 67,
                temperature: 6.7,
                relativeTemperature: 0.42,
                sharedWith: [{user: {name: 'Alice Johnson', email: 'alice@example.com'}}],
            },
        ];

        return (
            <AppBackground>
                <Box sx={{p: 6}}>
                    <StorySection
                        title="Combined States"
                        description="Default · Shared + Hot · Shared By · Long Name + Shared"
                    >
                        {albums.map((album) => (
                            <AlbumCard key={album.albumId.folderName} album={album} onClick={fn()} />
                        ))}
                    </StorySection>
                </Box>
            </AppBackground>
        );
    },
};

// ============================================================================
// 7. All States Overview (Comprehensive)
// ============================================================================

export const AllStates: Story = {
    render: () => (
        <AppBackground>
            <Box sx={{p: 6}}>
                <Typography
                    variant="h1"
                    sx={{
                        mb: 4,
                        fontSize: 28,
                        fontWeight: 300,
                        color: '#4a9ece',
                        textAlign: 'center',
                        letterSpacing: '0.08em',
                        textTransform: 'uppercase',
                    }}
                >
                    AlbumCard Component - All States
                </Typography>

                {/* Default */}
                <StorySection title="1. Default" description="Not shared, low temperature, typical photo count">
                    <AlbumCard
                        album={{
                            albumId: createAlbumId('owner-1', 'default'),
                            name: 'Family Vacation',
                            start: new Date('2026-07-15'),
                            end: new Date('2026-07-22'),
                            totalCount: 47,
                            temperature: 6.7,
                            relativeTemperature: 0.25,
                            sharedWith: [],
                        }}
                        onClick={fn()}
                    />
                </StorySection>

                {/* Shared */}
                <StorySection title="2. Shared" description="Shared with 1 person vs. Shared with 3 people">
                    <AlbumCard
                        album={{
                            albumId: createAlbumId('owner-1', 'shared-1'),
                            name: 'Summer Vacation 2026',
                            start: new Date('2026-08-01'),
                            end: new Date('2026-08-15'),
                            totalCount: 156,
                            temperature: 10.4,
                            relativeTemperature: 0.6,
                            sharedWith: [
                                {user: {name: 'Alice Johnson', email: 'alice@example.com', picture: 'https://i.pravatar.cc/150?img=1'}},
                            ],
                        }}
                        onClick={fn()}
                    />
                    <AlbumCard
                        album={{
                            albumId: createAlbumId('owner-1', 'shared-3'),
                            name: 'Paris Trip',
                            start: new Date('2026-03-10'),
                            end: new Date('2026-03-17'),
                            totalCount: 203,
                            temperature: 29.0,
                            relativeTemperature: 0.85,
                            sharedWith: [
                                {user: {name: 'Alice Johnson', email: 'alice@example.com', picture: 'https://i.pravatar.cc/150?img=1'}},
                                {user: {name: 'Bob Smith', email: 'bob@example.com', picture: 'https://i.pravatar.cc/150?img=2'}},
                                {user: {name: 'Carol White', email: 'carol@example.com'}},
                            ],
                        }}
                        onClick={fn()}
                    />
                </StorySection>

                {/* Temperature */}
                <StorySection title="3. Temperature Levels" description="From cold (0.1) to hot (1.0)">
                    {[
                        {name: 'Weekend Getaway', count: 12, temp: 0.1},
                        {name: 'City Tour', count: 34, temp: 0.35},
                        {name: 'Beach Holiday', count: 89, temp: 0.55},
                        {name: 'Iceland Adventure', count: 187, temp: 0.82},
                        {name: 'Christmas 2025', count: 234, temp: 1.0},
                    ].map((data, idx) => (
                        <AlbumCard
                            key={idx}
                            album={{
                                albumId: createAlbumId('owner-1', `temp-${idx}`),
                                name: data.name,
                                start: new Date('2026-01-01'),
                                end: new Date('2026-01-10'),
                                totalCount: data.count,
                                temperature: data.count / 10,
                                relativeTemperature: data.temp,
                                sharedWith: [],
                            }}
                            onClick={fn()}
                        />
                    ))}
                </StorySection>

                {/* Photo Count */}
                <StorySection title="4. Photo Count" description="0 to 3 photos">
                    {[0, 1, 2, 3].map((count) => (
                        <AlbumCard
                            key={count}
                            album={{
                                albumId: createAlbumId('owner-1', `count-${count}`),
                                name: `${count} Photo${count !== 1 ? 's' : ''}`,
                                start: new Date('2026-01-01'),
                                end: new Date('2026-01-01'),
                                totalCount: count,
                                temperature: count,
                                relativeTemperature: count * 0.05,
                                sharedWith: [],
                            }}
                            onClick={fn()}
                        />
                    ))}
                </StorySection>

                {/* Shared By */}
                <StorySection title="5. Shared By" description="Albums owned by other users">
                    <AlbumCard
                        album={{
                            albumId: createAlbumId('alice-owner', 'alice-vacation'),
                            name: "Alice's Summer Trip",
                            start: new Date('2026-06-15'),
                            end: new Date('2026-06-22'),
                            totalCount: 124,
                            temperature: 17.7,
                            relativeTemperature: 0.65,
                            ownedBy: {
                                name: 'Alice Johnson',
                                users: [{name: 'Alice Johnson', email: 'alice@example.com', picture: 'https://i.pravatar.cc/150?img=1'}],
                            },
                            sharedWith: [],
                        }}
                        onClick={fn()}
                    />
                    <AlbumCard
                        album={{
                            albumId: createAlbumId('family-owner', 'family-reunion'),
                            name: 'Family Reunion 2026',
                            start: new Date('2026-12-20'),
                            end: new Date('2026-12-27'),
                            totalCount: 287,
                            temperature: 35.9,
                            relativeTemperature: 0.92,
                            ownedBy: {
                                name: 'Family Account',
                                users: [
                                    {name: 'John Doe', email: 'john@example.com', picture: 'https://i.pravatar.cc/150?img=3'},
                                    {name: 'Jane Doe', email: 'jane@example.com', picture: 'https://i.pravatar.cc/150?img=4'},
                                ],
                            },
                            sharedWith: [],
                        }}
                        onClick={fn()}
                    />
                </StorySection>
            </Box>
        </AppBackground>
    ),
};
